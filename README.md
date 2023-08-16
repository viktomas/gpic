# gpic

Visual tool to quickly remove similar images from your photo folder. For people who take many almost-identical pictures and don't bother to clean up the "duplicates".

![showcase gif](./doc/gpic.gif)

## Usage:

- Download the binary from the release page or install `go install github.com/viktomas/gpic`
- `gpic folder/with/pictures`

From here, `gpic` opens a local web server that will show you the pictures and lets you remove some of the similar pictures.

**Don't worry, `gpic` never removes the images, it only moves them to a `./to-delete/` subfolder in your image folder.**
