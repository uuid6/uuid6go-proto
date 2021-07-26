# UUIDv6, 7 and 8

Welcome new contestors, UUIDv6, v7 and v8 in Golang. Nice, awesome, machine-sortable and absolutely unique!

 - UUIDv6 - a highly compatible guy.
 - UUIDv7 - all brand new and shinny. Very awesome.
 - UUIDv8 - strange and implementration specific, to fulfill everyone's dreams and requirements.

For more information visit https://datatracker.ietf.org/doc/draft-peabody-dispatch-new-uuid-format/

## Usage

```go
var gen uuid.UUIDv7Generator

//Sets how many bits of nano-second precision you would like to see in your UUID. Don't go under 12, don't exceed 48
gen.Precision = 12

id := gen.Next()


fmt.Println(id.ToString())
//Output:
//060f1cb1ce-8c70-e7b9-f7e4-4b50e3353e

fmt.Println(id.ToBinaryString())
//Output:
//00000110 00001111 00011100 10110001 11001110 10001100 01110000 11100111 10111001 11110111 11100100 01001011 01010000 11100011 00110101 00111110
fmt.Println(id.ToMicrosoftString())
//Output:
//{060F1CB1CE-8C70-E7B9-F7E4-4B50E3353E}

fmt.Println(id.Time())
//Ouput:
//2021-07-16 11:08:28 -0700 PDT

fmt.Println(id.Timestamp())
//Output:
//1626458908
```