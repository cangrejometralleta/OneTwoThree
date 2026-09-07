package main

import _ "embed"

// Liberation Serif Plays the Role Bitstream Charter Played in style.css.
// DejaVu Sans Mono Keeps the Mono Face style.css already Chose.
// Both Ship under free Licenses, See fonts/LIBERATION-LICENSE and
// fonts/DEJAVU-LICENSE.

//go:embed fonts/LiberationSerif-Regular.ttf
var fontSerifRegular []byte

//go:embed fonts/LiberationSerif-Bold.ttf
var fontSerifBold []byte

//go:embed fonts/LiberationSerif-Italic.ttf
var fontSerifItalic []byte

//go:embed fonts/LiberationSerif-BoldItalic.ttf
var fontSerifBoldItalic []byte

//go:embed fonts/DejaVuSansMono.ttf
var fontMonoRegular []byte
