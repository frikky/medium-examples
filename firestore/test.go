package user

type User struct {
	Username   string   `json:"userId"`
	Email      string   `json:"email"`
	DateEdited int64    `json:"date_edited"`
	Access     []Access `json:"access"`
}

type Access struct {
	Id   string `json:"id"`
}

user := User{
	Username: "username",
	Access: []Access{
		Id: "id",
	},
}

type User struct {
	Username   StringValue  `json:"userId"`
	Email      StringValue  `json:"email"`
	DateEdited IntegerValue `json:"date_edited"`
	Access     []Access     `json:"access"`
}

type AccessWrapper struct {
	ArrayValue ArrayValue `json:"arrayValue"`
}

type Access struct {
	Id   StringValue `json:"id"`
}

type ArrayValue struct {
	Values []Value `json:"values"`
}

type Value struct {
	MapValue     MapValue     `json:"mapValue,omitempty"`
	StringValue  StringValue  `json:"stringValue,omitempty"`
	IntegerValue IntegerValue `json:"integerValue,omitempty"`
	ArrayValue   ArrayValue   `json:"arrayValue,omitempty"`
}

type IntegerValue struct {
	IntegerValue string `json:"integerValue"`
}

type StringValue struct {
	StringValue string `json:"stringValue"`
}

user := User{
	Username: StringValue{
		StringValue: "username",
	},
	AccessWrapper: AccessWrapper{
		ArrayValue: ArrayValue{
			Values: []Value{
				Value{
					MapValue: MapValue{
						Access: Access{
							Id: StringValue{
								StringValue: "id",
							},
						},
					},
				},
			},
		},
	}
}
