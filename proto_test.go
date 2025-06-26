package main

// // Don't need it anymore
// func parseCommand(raw string) (Command, error) {
// 	rd := resp.NewReader(bytes.NewBufferString(raw))
// 	v, _, err := rd.ReadValue()

// 	if err == io.EOF {
// 		return nil, err
// 	}

// 	if v.Type() == resp.Array {
// 		for _, value := range v.Array() {
// 			switch value.String() {
// 			case CommandSet:
// 				if len(v.Array()) != 3 {
// 					return nil, fmt.Errorf("invalid number of variables for SET command")
// 				}
// 				cmd := SetCommand{
// 					key: v.Array()[1].Bytes(),
// 					val: v.Array()[2].Bytes(),
// 				}
// 				return cmd, nil
// 			case CommandGet:
// 				if len(v.Array()) != 2 {
// 					return nil, fmt.Errorf("invalid number of variables for GET command")
// 				}
// 				cmd := GetCommand{
// 					key: v.Array()[1].Bytes(),
// 				}
// 				return cmd, nil
// 			case CommandHello:
// 				cmd := HelloCommand{
// 					value: v.Array()[1].String(),
// 				}
// 				return cmd, nil
// 			}
// 		}
// 	}

// 	return nil, fmt.Errorf("invalid or unknown command recieved : %s", raw)
// }

// func TestProtocol(t *testing.T) {
// 	raw := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
// 	cmd, err := parseCommand(raw)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println(cmd)
// }
