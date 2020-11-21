package common

type UINT64 uint64

//type INT64 int64
const TrafficKUnit int = 1000.0

// PricePerByte = PricePerGB/(TrafficKUnit*TrafficKUnit*TrafficKUnit)
// in DB region_price 10/KB=>about $0.01/GB (depends on TrafficKUnit)
const RateToOneDollar uint64 = 1e9 // amount / 1e9 = real Dollar
