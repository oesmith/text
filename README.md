# Text

Some text algorithms in Go.

## Usage

    import "github.com/oesmith/text"

### Double Metaphone

    pri, sec = text.DoubleMetaphone("Jackson")
    // => "JKSN", "AKSN"
    pri, sec = text.DoubleMetaphone("Rowland")
    // => "RLNT, "" (no secondary)

## Thanks

* @threedaymonk, for the Ruby algorithms from which this code was translated
