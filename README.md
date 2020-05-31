# REST & GRPC SERVER

## How to run
- Generate protobuf
<br>
`cd pb && protoc movie.proto --go_out=plugins=grpc:.`


- Instal depedency
<br>
`go mod download`

- Run 
<br>
`go run main.go`






# Query
```
SELECT u1.ID, U1.Username, u2.UserName as Parent
    FROM USER u1 LEFT JOIN USER u2
        ON u1.Parent = u2.ID
```