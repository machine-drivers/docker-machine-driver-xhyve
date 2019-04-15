// +build qcow2

package hyperkit

//go:generate ocamlfind ocamlopt -thread -package "cstruct cstruct.lwt io-page io-page.unix uri mirage-block mirage-block-unix qcow unix threads lwt lwt.unix" -c mirage_block_ocaml.ml -o mirage_block_ocaml.cmx
//go:generate ocamlfind ocamlopt -thread -linkpkg -package "cstruct cstruct.lwt io-page io-page.unix uri mirage-block mirage-block-unix qcow unix threads lwt lwt.unix" -output-obj -o mirage_block_ocaml.o mirage_block_ocaml.cmx

/*
// Not need ${SRCDIR} because '-l' linker flag should not include the file path
#cgo LDFLAGS: -lmirage_block_ocaml.o
*/
import "C"
