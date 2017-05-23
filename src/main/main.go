package main

func main()  {
	a := App{}
	a.Initialize("postgres://mjaovdiv:xxiD03VhsirEPIc-bX171XXchv_3vQhb@horton.elephantsql.com:5432/mjaovdiv")
	a.Run(":3333")
}

