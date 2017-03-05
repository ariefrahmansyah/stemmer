stemmer
==========

Stemmer for Indonesian Language

    import "github.com/ariefrahmansyah/stemmer"

Installation

    go get github.com/ariefrahmansyah/stemmer

Usage:

    stemm := stemmer.New()
    out := stemm.Stemm("memakan")
    println(out[0])
    
    out = stemm.Stemm("memakan", "mencintai")
    println(out[0], out[1])