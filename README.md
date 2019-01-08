# grapheme-splitter

# 背景

尽管Golang使用utf8作为默认编码，但是在处理某些特殊字符时仍然不能准确判断字符串的长度

例如，普通emoji字符 "🌷","🎁","💩","😜" 和 "👍" 可以确定长度为1，但是像国旗"🇨🇳"这样复杂且常用的emoji的长度就不准确了 

```golang
len([]rune("🌷")) == 1

len([]rune("🇨🇳")) == 2
```

一些不常用的汉字，或者其他国家的特殊文字识别出的长度与肉眼所见会存在误差
    
```golang
two := "ñ"; // unnormalized two-char n+◌̃  , i.e. "\u006E\u0303";
one := "ñ"; // normalized single-char, i.e. "\u00F1"
len([]rune(two)) // 2
len([]rune(one)) // 1
```

本库的目的在于将任意的字符串切割成肉眼所见的单个字符

# 安装

通过如下命令进行安装:

```bash
$ go get github.com/SubLuLu/grapheme-splitter
```

# 测试

To run the tests on `grapheme-splitter`, use the command below:

```bash
$ go test -v

$ go test -bench=.
```

# 使用

`import`成功后即可使用:

```goalng
import splitter "github.com/SubLuLu/grapheme-splitter"

// split the string to an array of grapheme clusters (one string each)
graphemes := splitter.Split(string)

// or do this if you just need their number
counter := splitter.Counter(string)
```

# 示例

```golang
import splitter "github.com/SubLuLu/grapheme-splitter"

// plain latin alphabet - nothing spectacular
splitter.Split("abcd"); // returns ["a", "b", "c", "d"]

// two-char emojis and six-char combined emoji
splitter.Split("🌷🎁💩😜👍🇨🇳"); // returns ["🌷","🎁","💩","😜","👍","🇨🇳"]
```
