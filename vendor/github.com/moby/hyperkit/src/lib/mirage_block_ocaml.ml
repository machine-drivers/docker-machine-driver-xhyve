(* Initialise logging *)
let src =
  Logs.set_reporter (Logs_fmt.reporter ());
  let src = Logs.Src.create "mirage" ~doc:"mirage block device interface" in
  Logs.Src.set_level src (Some Logs.Info);
  src
module Log = (val Logs.src_log src : Logs.LOG)


(* TODO: parameterise this over block implementations *)
module Qcow = Qcow.Make(Block)(OS.Time)

module Mutex = struct
  include Mutex

  let with_lock m f =
    Mutex.lock m;
    try
      let r = f () in
      Mutex.unlock m;
      r
    with
    | e ->
      Mutex.unlock m;
      raise e
end

(* The Lwt thread will signal the pthread thread by filling an Ivar, which
   briefly locks a mutex and signals on a regular condition variable. This
   is considered non-blocking and hence safe to call from Lwt. *)
module Ivar = struct
  type 'a t = {
    mutable item: 'a option;
    m: Mutex.t;
    c: Condition.t;
  }

  let make () = { item = None; m = Mutex.create (); c = Condition.create () }

  let wait t =
    Mutex.with_lock t.m
      (fun () ->
        let rec loop () = match t.item with
          | None ->
            Condition.wait t.c t.m;
            loop ()
          | Some x ->
            x in
        let result = loop () in
        result
      )

  let fill t x =
    Mutex.with_lock t.m
      (fun () ->
        if t.item <> None then begin
          Printf.fprintf stderr "Ivar already filled\n%!";
          exit 1
        end;
        t.item <- Some x;
        Condition.signal t.c;
      )
end


module Protocol = struct
  module Request = struct
    type t =
      | Connect of (Block.Config.t * Qcow.Config.t option)
      | Get_info of int
      | Disconnect of int
      | Read of int * int * (Cstruct.t list)
      | Write of int * int * (Cstruct.t list)
      | Delete of int * int64 * int64
      | Flush of int
  end

  module Response = struct
    type ok =
      | Connect of int
      | Get_info of bool * int * int64 * bool
      | Disconnect
      | Read of int
      | Write of int
      | Delete
      | Flush
    type t = (ok, Qcow.write_error) result
  end

  (* An in-flight request *)
  type t = {
    request: Request.t;
    ivar: Response.t Ivar.t;
  }

  (* The pthread code signals the Lwt code by writing a byte to the pipe.
     In an ideal world we would use an Lwt.task but the implementation of
     wakeup is not pthread-safe (in particular there are unguarded Queue.push
     calls) *)
  let request_reader, request_writer = Unix.pipe()

  (* The list of in-flight requests is shared by both client and server.
     The mutex is held (briefly) to ensure the integrity of the list. *)
  let in_flight_requests = ref []
  let in_flight_requests_m = Mutex.create ()

  (* Called by the server to retrieve all the in-flight requests.
     Considered non-blocking since the mutex is only ever held to
     manipulate the list. *)
  let take_all () =
    Mutex.with_lock in_flight_requests_m
      (fun () ->
        let results = !in_flight_requests in
        in_flight_requests := [];
        List.rev results
      )

  (* Called by the client to send a request to the server. This will
     block the calling pthread until the signal has been sent. *)
  let send request =
    let ivar = Ivar.make () in
    let t = { request; ivar } in
    Mutex.with_lock in_flight_requests_m
      (fun () ->
        in_flight_requests := t :: !in_flight_requests;
      );
    let n = Unix.write request_writer "X" 0 1 in
    if n = 0 then begin
      Printf.fprintf stderr "Got EOF while writing signal to the pipe\n%!";
      exit 1;
    end;
    t

  (* Called by the client to perform an RPC. This will block the calling
     pthread until the result has been created. *)
  let rpc request : Response.t =
    let t = send request in
    Ivar.wait t.ivar

end

module C = struct
  (** The C callbacks live here. These are all called from pthreads and are
      allowed to block. *)

  let string_of_error = function
    | `Unimplemented -> "Operation is not implemented"
    | `Is_read_only -> "Block device is read-only"
    | `Disconnected -> "Block device is disconnected"
    | _ -> "Unknown error"

  let ok_exn = function
    | Error x ->
      Printf.fprintf stderr "Mirage-block error: %s\n%!" (string_of_error x);
      failwith (string_of_error x)
    | Ok x -> x

  let mirage_block_open block_config qcow_config_opt stats_config_opt : int =
    let description = Printf.sprintf "block_config = %s and qcow_config = %s and stats_config = %s"
      block_config
      (match qcow_config_opt with None -> "None" | Some x -> x)
      (match stats_config_opt with None -> "None" | Some x -> x)
      in
    Printf.fprintf stdout "mirage_block_open: %s\n%!" description;
    match Block.Config.of_string block_config with
      | Ok block_config' ->
        let qcow_config' = match qcow_config_opt with
          | None -> None
          | Some x ->
            begin match Qcow.Config.of_string x with
              | Ok qcow_config' -> Some qcow_config'
              | Error (`Msg m) ->
                Printf.fprintf stderr "mirage_block_option %s: %s\n%!" description m;
                exit 1
            end in
        (* Start exposing stats if configured *)
        begin match stats_config_opt with
          | None -> ()
          | Some x ->
            let mode =
              (* either unix:<path> or tcp:<port> *)
              match Astring.String.cut ~sep:":" x with
              | Some ("unix", path) -> `Unix_domain_socket (`File path)
              | Some ("tcp", port) -> `TCP (`Port (int_of_string port))
              | _ ->
                Printf.fprintf stderr "mirage_block_open: unrecognised stats_config: %s" x;
                exit 1 in
            let module Server = Prometheus_app.Cohttp(Cohttp_lwt_unix.Server) in
            let callback = Server.callback in
            let server = Cohttp_lwt_unix.Server.make ~callback () in
            Lwt.async (fun () -> Cohttp_lwt_unix.Server.create ~mode server)
        end;
        begin match ok_exn (Protocol.rpc (Protocol.Request.Connect (block_config', qcow_config'))) with
          | Protocol.Response.Connect t ->
            Printf.fprintf stdout "mirage_block_open: %s returning %d\n%!" description t;
            t
          | _ ->
            Printf.fprintf stderr "protocol error: unexpected response to connect\n%!";
            exit 1
        end
      | Error (`Msg m) ->
        Printf.fprintf stderr "mirage_block_open %s: %s\n%!" description m;
        exit 1

  let mirage_block_stat (h: int) : (bool * int * int64 * bool) =
    Printf.fprintf stdout "mirage_block_stat\n%!";
    match ok_exn (Protocol.rpc (Protocol.Request.Get_info h)) with
      | Protocol.Response.Get_info (read_write, sector_size, size_sectors, candelete) ->
        read_write, sector_size, size_sectors, candelete
      | _ ->
        Printf.fprintf stderr "protocol error: unexpected response to stat\n%!";
        exit 1

  let mirage_block_preadv (h: int) (bufs: Io_page.t array) (ofs: int) : int =
    let bufs = Array.(to_list (map Io_page.to_cstruct bufs)) in
    match ok_exn (Protocol.rpc (Protocol.Request.Read (h, ofs, bufs))) with
      | Protocol.Response.Read len ->
        len
      | _ ->
        Printf.fprintf stderr "protocol error: unexpected response to read\n%!";
        exit 1

  let mirage_block_pwritev (h: int) (bufs: Io_page.t array) (ofs: int) : int =
    let bufs = Array.(to_list (map Io_page.to_cstruct bufs)) in
    match ok_exn (Protocol.rpc (Protocol.Request.Write (h, ofs, bufs))) with
      | Protocol.Response.Write len ->
        len
      | _ ->
        Printf.fprintf stderr "protocol error: unexpected response to write\n%!";
        exit 1

  let mirage_block_delete (h: int) (ofs: int64) (len: int64) : unit =
    match ok_exn (Protocol.rpc (Protocol.Request.Delete(h, ofs, len))) with
      | Protocol.Response.Delete ->
        ()
      | _ ->
        Printf.fprintf stderr "protocol error: unexpected response to delete\n%!";
        exit 1

  let mirage_block_flush h =
    match ok_exn (Protocol.rpc (Protocol.Request.Flush h)) with
      | Protocol.Response.Flush ->
        ()
      | _ ->
        Printf.fprintf stderr "protocol error: unexpected response to flush\n%!";
        exit 1

  let mirage_block_close h =
    match ok_exn (Protocol.rpc (Protocol.Request.Disconnect h)) with
      | Protocol.Response.Disconnect ->
        ()
      | _ ->
        Printf.fprintf stderr "protocol error: unexpected response to disconnect\n%!";
        exit 1

  let () =
    Callback.register "mirage_block_open" mirage_block_open;
    Callback.register "mirage_block_stat" mirage_block_stat;
    Callback.register "mirage_block_close" mirage_block_close;
    Callback.register "mirage_block_preadv" mirage_block_preadv;
    Callback.register "mirage_block_pwritev" mirage_block_pwritev;
    Callback.register "mirage_block_delete" mirage_block_delete;
    Callback.register "mirage_block_flush" mirage_block_flush;
end

(* Allocate connected block devices integers, to pass back to C as
   "file descriptors" *)
module Handle = struct

  type t = {
    base: Block.t;
    block: Qcow.t;
    info: Mirage_block.info;
  }

  let table = Hashtbl.create 7

  let register : t -> int =
    let next_id = ref 0 in
    fun t ->
      let this_id = !next_id in
      incr next_id;
      Hashtbl.replace table this_id t;
      this_id

  let deregister id = Hashtbl.remove table id

  let find_or_quit h =
    if not (Hashtbl.mem table h) then begin
      Printf.fprintf stderr "FATAL: mirage_block_wrapper failed to find open handle %d\n%!" h;
      exit 1
    end;
    Hashtbl.find table h
end

(* Below here is all the Lwt code *)

open Lwt

(* Process one request, send the response back to the pthread thread by filling
   the Ivar. *)
let process_one t =
  let open Protocol in
  let module Monad = struct
    let (>>=) m f =
      let open Lwt.Infix in
      m >>= function
      | Error e -> Lwt.return (Error e)
      | Ok x -> f x
    let return x = Lwt.return (Ok x)
  end in
  let open Monad in
  let result_t =
    Lwt.catch
      (fun () ->
        match t.request with
          | Request.Connect (block_config, qcow_config) ->
            let open Lwt.Infix in
            Block.of_config block_config
            >>= fun base ->
            Qcow.connect ?config:qcow_config base
            >>= fun block ->
            Qcow.get_info block
            >>= fun info ->
            let h = Handle.register { Handle.block; base; info } in
            return (Response.Connect h)
          | Request.Get_info h ->
            let t = Handle.find_or_quit h in
            let open Lwt.Infix in
            Qcow.get_info t.Handle.block
            >>= fun info ->
            let config = Qcow.to_config t.Handle.block in
            let candelete = config.Qcow.Config.discard in
            return (Response.Get_info (info.Mirage_block.read_write, info.Mirage_block.sector_size, info.Mirage_block.size_sectors, candelete))
          | Request.Disconnect h ->
            let t = Handle.find_or_quit h in
            let open Lwt.Infix in
            Qcow.disconnect t.Handle.block
            >>= fun () ->
            return Response.Disconnect
          | Request.Read (h, offset, bufs) ->
            let t = Handle.find_or_quit h in
            (* Offset needs to be translated into sectors *)
            if offset mod t.Handle.info.Mirage_block.sector_size <> 0 then begin
              Printf.fprintf stderr "Read offset not at sector boundary\n%!";
              exit 1
            end;
            let sector = Int64.of_int (offset / t.Handle.info.Mirage_block.sector_size) in
            let open Lwt.Infix in
            (* Qcow.error <> Qcow.write_error *)
            begin Qcow.read t.Handle.block sector bufs
              >>= function
              | Error `Disconnected -> Lwt.return (Error `Disconnected)
              | Error `Unimplemented -> Lwt.return (Error `Unimplemented)
              | Ok () ->
              let len = List.(fold_left (+) 0 (map Cstruct.len bufs)) in
              return (Response.Read len)
            end
          | Request.Write (h, offset, bufs) ->
            let t = Handle.find_or_quit h in
            (* Offset needs to be translated into sectors *)
            if offset mod t.Handle.info.Mirage_block.sector_size <> 0 then begin
              Printf.fprintf stderr "Write offset not at sector boundary\n%!";
              exit 1
            end;
            let sector = Int64.of_int (offset / t.Handle.info.Mirage_block.sector_size) in
            Qcow.write t.Handle.block sector bufs
            >>= fun () ->
            let len = List.(fold_left (+) 0 (map Cstruct.len bufs)) in
            return (Response.Write len)
          | Request.Delete (h, offset, len) ->
            let t = Handle.find_or_quit h in
            (* Offset and len need to be translated into sectors *)
            let sector_size = Int64.of_int t.Handle.info.Mirage_block.sector_size in
            if Int64.rem offset sector_size <> 0L then begin
              Printf.fprintf stderr "Delete offset not at sector boundary\n%!";
              exit 1
            end;
            if Int64.rem len sector_size <> 0L then begin
              Printf.fprintf stderr "Delete len not a multiple of sectors\n%!";
              exit 1
            end;
            let sector = Int64.div offset sector_size in
            let n = Int64.div len sector_size in
            Qcow.discard t.Handle.block ~sector ~n ()
            >>= fun () ->
            return Response.Delete
          | Request.Flush h ->
            let t = Handle.find_or_quit h in
            let open Lwt.Infix in
            (* Block.write_error <> Qcow.write_error *)
            begin Block.flush t.Handle.base
              >>= function
              | Error `Disconnected -> Lwt.return (Error `Disconnected)
              | Error `Is_read_only -> Lwt.return (Error `Is_read_only)
              | Error `Unimplemented -> Lwt.return (Error `Unimplemented)
              | Ok () ->
                return Response.Flush
            end
      )
      (fun e ->
        (* If the thread fails unexpectedly then convert this into an Error return
           to guarantee we return an error code to the block interface instead of
           dropping the request on the floor. *)
        Log.err (fun f -> f "Mirage block device raised exception: %s" (Printexc.to_string e));
        Lwt.return (Error `Disconnected)
      ) in
  let open Lwt in
  result_t >>= fun result ->
  Ivar.fill t.ivar result;
  return ()

(* An Lwt thread which receives the signals from the pipe, grabs the in-flight
   requests and forks background threads to process all the requests. *)
let serve_forever () =
  let buf = String.make 1 '\000' in
  let request_reader = Lwt_unix.of_unix_file_descr Protocol.request_reader in

  let rec loop () =
    Lwt_unix.read request_reader buf 0 1
    >>= fun n ->
    if n = 0 then begin
      Printf.fprintf stderr "Got EOF while reading signal from the pipe\n%!";
      exit 1
    end;
    let all = Protocol.take_all () in
    let (_: unit Lwt.t list) = List.map process_one all in
    loop () in
  loop ()
let (_: Thread.t) = Thread.create (fun () -> Lwt_main.run (serve_forever ())) ()
