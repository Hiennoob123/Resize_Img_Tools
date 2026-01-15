# RESIZE IMAGE TOOL

## Goals
- Tools for training image detection model
- Resizing images to lower size or resolution -> Required less memory to store
- Utilize Multithread to handle jobs -> Improve efficiency

## Prerequisite
- Compilers: Go, C
- Libvips: https://github.com/libvips/libvips 

Or, you can install all the tools by running

```bash
chmod +x install_tool.sh
sudo ./install_tool.sh
```

## Usage
- Operations: Resize Image to 720p or 1080p
- Flags:
--r: Recursively walk through the source directory and maintain the same relative path for output
--res: Specify Resolution (p720p or p1080)

- Examples command:

```bash
go run . <src_dir> <dst_dir>
go run . --r <src_dir> <dst_dir>
go run . --res 1080p --r <src_dir> <dst_dir>
```



## Credit
This project uses govips (MIT License):
https://github.com/davidbyttow/govips


