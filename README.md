# 几个简单的golang工具

## desc 
### 从终端copy如下,可以转换成

```
| id    | int(11)  | YES  |     | NULL    |       |
| dt    | datetime | YES  |     | NULL    |       |
+-------+----------+------+-----+---------+-------+
```

```
"id","dt"
```

### 可以比较两个字段的diff

```
| id    | int(11)  | YES  |     | NULL    |       |
| dt    | datetime | YES  |     | NULL    |       |
+-------+----------+------+-----+---------+-------+
```

```
| id    | int(11)  | YES  |     | NULL    |       |
| dt    | datetime | YES  |     | NULL    |       |
| dt2    | datetime | YES  |     | NULL    |       |
+-------+----------+------+-----+---------+-------+
```

输出不同点(仅仅是字段,不包括字段的属性):

```
Len1: 2  Len2: 3
Diff: [dt2]
```

## insert
### 可以吧下面的这种,转化成insert语句

```
|    1 | 2019-03-05 01:53:55 |
|    1 | 12                |
```

## excel

### 在两个excel之间操作,经行简单的运算
```
type 操作类型错误 opType:
Usage of ./main:
  -a string
    	a excel文件
  -acol string
    	a excel文件需要比较的列的头名称,多个以逗号隔开
  -b string
    	b excel文件
  -bcol string
    	b excel文件需要比较的列的头名称,多个以逗号隔开
  -debug
    	是否debug debug输出一些中间日志
  -output-file string
    	输出类型文件名称 为空只输出到屏幕
  -sheet string
    	excel文件工作表名称 (default "Sheet1")
  -type string
    	操作类型 [a-b] a有的b没有 [b-a] b有的a没有 [a+b] a和b的并集(去重) [ab] a和b的都有的,即交集
```		