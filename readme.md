# findImagePosition  [中文文档](https://github.com/topascend/findImagePosition/edit/main/readme_zh.md)
## Search for the top left corner coordinates of the small image in the large image
##
```
go get github.com/topascend/findImagePosition
```
##
#### FindPosition 
- returns the position of the small image in the large image, and returns false if not found
####  FindAnyPosition 
- single matching version, stop searching if any position is found, return false if not found
#### FindAllPositions 
- returns all matching positions, false if not found
### Testing

![大图](./examples/big.png)
![小图](./examples/small0.png)
![小图](./examples/small4.png)



```
findPosition taked time:  10.2394ms
Found at: 8 10
findAnyPosition taked time:  10.0666ms
Found at: 8 10
findAllPositions taked time:  15.9491ms
[(8,10)]
Found at: [(8,10)]
picture: examples\small0.png taked total time: : 54.7038ms


findPosition taked time:  7.0034ms
Found at: 1680 37
findAnyPosition taked time:  6.982ms
Found at: 1680 37
findAllPositions taked time:  7.975ms
[(1680,37) (1682,141) (1684,253) (1684,367) (1686,843)]
Found at: [(1680,37) (1682,141) (1684,253) (1684,367) (1686,843)]
picture: examples\small4.png taked total time: : 36.5171ms 
```
