# img2sphere

Thanks for visiting, this is a golang utility for working with recti-linear textures and images that we would like to project stereographically with a conformal mapping.

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
2. Quadratic - $x^2$ projection [-q]
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


