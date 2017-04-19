using Go = import "/go.capnp";
@0xf4533cbae6e08506;
$Go.package("main");
$Go.import("main");

struct TlogBlock {
	sequence @0 :UInt64;
	text @1 :Text;
}

struct TlogAggregation {
	name @0 :Text; # unused now
	size @1 :UInt64; # number of blocks in this aggregation
	timestamp @2 :UInt64;
	volumeId @3 :UInt32;
	blocks @4 :List(TlogBlock);
	prev @5 :Data; # hash of the previous aggregation
}
