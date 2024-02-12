# img2sphere

Thanks for visiting, this is a golang utility for working with recti-linear textures and images that we would like to project stereographically with a conformal mapping.

The tool works with a process in unwrapping in blender where we would take the following steps:

1. Add a UV Sphere
2. Enter Edit Mode
3. Alt + Left Click to Select Edge Loop Around the Equator
4. Edges -> Mark Seam
5. Enter UV Editor
6. Ctrl + A Select all Sphere Points
7. U Unwrap Vertices (Should give us a hemisphere)

By apply default mode conformal projection (linear) to a rectilinear texture we can map it to this sphere without any distortion. Note that on a disk:

$$ dA = dr^2 \theta $$


## installation

```
go install github.com/andewx/img2sphere
```

Installs `./img2sphere` in your `/usr/local/bin` `$GOBIN` path

## Usage



Once the tool is installed you can use this tool to map rectilinear texturest to hemispheric conformal projections.

```
./img2sphere [-h] <input_file> <output_file> <options>
```

Our tool supports three mappings which refer to the equatorial spherical distortion:

1. Linear Magnitude projection - where the equatorial distorition is linear [default]
2. Quadratic - $x^2$ projection [-q]. This will stretch the texture by a factor along its mapped radial
3. Cubic - $x^3$ projection [-c]
4. Log - natural log projection [-ln]
5. Exp - Exponential projection. [-x]

To use simply use any of these options.

## Options

-h - Shows help
-s - Set scale
-l - Log uv mapping transforms
-c - cubic mode
-q - quadratic mode
-ln - log mode
-x - exponential mode


## Support

This project likely will just be issued as is but if there are any suggestions or issues I will gladly handle them.

```
img2sphere inputfile.png example.png -ln
```

![Example](example.png)


