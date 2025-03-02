# findImagePosition
## Search for the top left corner coordinates of the small image in the large image
##
#### FindPosition 
- returns the position of the small image in the large image, and returns false if not found
####  FindAnyPosition 
- single matching version, stop searching if any position is found, return false if not found
#### FindAllPositions 
- returns all matching positions, false if not found
### Testing
```
FindPosition takes 10.4379ms
Found at: 1680 37
FindAnyPosition takes 9.2486ms
Found at: 1680 37
FindAllPositions takes 9.7378ms
[(1680,37) (1682,141) (1684,253) (1684,367) (1686,843)]
Found at: [(1680,37) (1682,141) (1684,253) (1684,367) (1686,843)]
Image: Example \ small4.png Total time: 49.8152ms 
```