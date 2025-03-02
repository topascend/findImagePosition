# findImagePosition

## 在大图中查找小图的左上角坐标

##
#### FindPosition 
- 返回小图在大图中的位置，若未找到返回false
#### FindAnyPosition 
- 单匹配版本,找到任意一个位置就停止查找，若未找到返回false
#### FindAllPositions 
- 返回所有匹配位置，若未找到返回false

### 测试

```
findPosition 用时 10.4379ms
Found at: 1680 37
findAnyPosition 用时 9.2486ms
Found at: 1680 37
findAllPositions 用时 9.7378ms
[(1680,37) (1682,141) (1684,253) (1684,367) (1686,843)]
Found at: [(1680,37) (1682,141) (1684,253) (1684,367) (1686,843)]
图片: examples\small4.png 总共用时: 49.8152ms 
```
