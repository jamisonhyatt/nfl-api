package nfl


func  Right (s string, i int) string {
    len := len(s)
    if len <= i {
        return s
    }
    s= string(s[len-i:len])
    return s;
}