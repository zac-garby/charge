# charge

[(read the blog post for a more detailed explaination)](http://zacgarby.co.uk/posts/1.html)

Suppose you have a cuboid containing some charged particles in random positions and with random (positive or negative) charges. For each slice of a certain height in that cuboid, look at each point on that plane and measure the force due to electric charge and convert the magnitude of that force using some function into a colour, then set the corresponding pixel of the corresponding frame of a gif to that colour. For each pixel, do the same, and then look at the gif. This program does exactly that, and it looks very cool.

![](out.gif)
 > Each frame of the gif is in the subdirectory 'out' as a png file.

More examples of this are available on [my website](https://zacgarby.co.uk).

## Usage

Download the program (`go get github.com/zac-garby/charge`) and also [install gifgen](https://github.com/lukechilds/gifgen).

Run the commands `rm -r out && mkdir out && go run main.go && gifgen -o out.gif -f 20 out/%d.png` to generate the output in the same way I am doing it -- remove all previous output frames, make the output directory, run the program (generates the frames), generate a 20fps gif from the images in the 'out' directory.
