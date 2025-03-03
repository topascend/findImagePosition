# findImagePosition

## 在大图中查找小图的左上角坐标
##
```
go get github.com/topascend/findImagePosition
```
##
#### FindPosition 
- 返回小图在大图中的位置，若未找到返回false
#### FindAnyPosition 
- 单匹配版本,找到任意一个位置就停止查找，若未找到返回false
#### FindAllPositions 
- 返回所有匹配位置，若未找到返回false

### 测试

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
