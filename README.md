# mvpls
Mvpls (move please) is a simple utility tool written in Go to allow for the recursive moving, copying, or removing of files and flattening of file trees with regex.

![](assets/mvpls.gif)

In this example, every pdf file within my books directory and its subdirectories is moved to a seperate folder.

Why are my pdfs executable? Idk calibre did that to them for some reason.

It's really simple, pass in a -r "REGEX" flag and the directories you would like to search within.
```sh
mvpls -r ".*\.png" . moveLocation/
```
moves all pngs extensions from the current directory and its subdirectories to the folder `moveLocation`.

to copy, pass in the -c flag, remove pass in the --remove flag.

it's still in development lol. Try at your own risk I suppose. Should work though... spent a few hours on it...

You could say, can't we easily do this by piping a few [standard commands](https://superuser.com/a/1041895/1225558)?
And you'd be right. But I wanted to make this anyway :P. Also  Go is fun to code in so I can't be wrong there.

It also works just like the POSIX `mv` command so can move files like regular when no flag is passed. Do I get a POSIX compliant badge??

## Why GO?
The correct question is why not Go? I mostly chose Go because it looks fun and is cross platform. 
The two most import parts in software development after the name.

### Pressing features
- Key feature allowing you to thank mvpls if it succeeds.
- Cache operations to allow for a possible revert flag
- Obviously add the ability to copy or remove.
