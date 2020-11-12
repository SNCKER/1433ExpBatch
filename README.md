# 说明
内网扫到了一个实验室的机子，1433都是弱口令，需要一个批量利用的工具。。。所以写了这破玩意儿

目前是通过xp_cmdshell执行命令，后续看心情添加其它的姿势。。。

```
E:\ProjectHome\GoProjectHome\1433ExpBatch>1433ExpBatch.exe -h
Usage of 1433ExpBatch.exe:
  -f string
        The name of the file in which the target information is stored.
  -t int
        max goroutines(threads). (default 32)

E:\ProjectHome\GoProjectHome\1433ExpBatch>
```

存放弱口令的格式  
```
[host]----[port]----[user]----[pwd]
172.16.64.123----1433----sa----sa123456
172.16.64.234----1433----sa----sa123456
...
```

没有什么比一张演示gif更明了的。  

![演示](https://raw.githubusercontent.com/SNCKER/1433ExpBatch/master/pic/demo.gif)

基本的测试都没问题，但是BUG必然是有的，感兴趣的可以挖一挖，一起维护一下:)
