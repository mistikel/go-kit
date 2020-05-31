SELECT u1.ID, U1.Username, u2.UserName as Parent
    FROM USER u1 LEFT JOIN USER u2
        ON u1.Parent = u2.ID