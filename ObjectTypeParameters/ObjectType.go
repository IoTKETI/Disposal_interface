package ObjectTypeParameters


//============================

type Mif struct {
	Mif []MifObject `json:"mif"`
}

type MifObject struct {
	Mi int  `json:"mi"`
	Ops []string `json:"ops"`
}


//============================ 들어오는 InPut에 따라서 달라짐

type Ip struct {
	Ip IpObject `json:"ip"`
}

type IpObject struct {
	Value int `json:"value"`
}

//============================ 들어오는 InPut에 따라서 달라짐

type Op struct {
	Op OpObject `json:"op"`
}

type OpObject struct {
	Value int `json:"value"`
}


//============================

type Tif struct {
	Tif []TifObject `json:"tif"`
}

type TifObject struct {   /// 이걸 만들어야함
	Ti int `json:"ti"`
	Sp string `json:"sp,omitempty"`  //들어올때마다 달라짐
	Fp string `json:"fp,omitempty"`
	To bool `json:"to,omitempty"`
}

//============================ 들어오는 InPut에 따라서 달라짐
//
type Fp struct {
	Oprd string `json:"oprd,omitempty"`
	Sprd string `json:"sprd,omitempty"`
}






//type FpObject struct {
//	Value int `json:"value"`
//}

//============================ 들어오는 InPut에 따라서 달라짐
//
//type Sp struct {
//	Sp string
//}

//type SpObject struct {
////	Value int `json:"value"`
//	Value int `json:"mast"`
//
//}

