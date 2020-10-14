package data

import (
	"sort"

	"github.com/findy-network/findy-agent-api/graph/model"
)

func init() {
	sort.Slice(connections, func(i, j int) bool {
		return connections[i].Created() < connections[j].Created()
	})

	sort.Slice(events, func(i, j int) bool {
		return events[i].Created() < events[j].Created()
	})
}

var user = InternalUser{
	"035f09fe-7e21-4546-a427-d954251f7082", "Carissa", 1453063800,
}

var connections = []InternalPairwise{
	{"95b2eab7-d664-4868-b684-3dbbcc3e0375", "ETbgvUngPJyPMIDJbDmoAogAS", "sxxyLTYTIsipWVqcXFPGwUahc", "https://www.hEPvSuS.net/lbOBXiP", "Greenholt Agency", true, 1183630175, 706882240},
	{"0a7f8386-833d-4c2c-9ffd-daf0805242f2", "DtrbPMCagUXDlaFddnsewtcgs", "cQXRRuWuHmShTHXolcqWSOKYk", "http://sZPgmbo.info/", "Wiegand Ltd", true, 1277848304, 1108656723},
	{"bd5b6b66-a2cd-451f-b906-37ea3dbb1301", "QdECjfqjnKknVkOhGkyHTWWfF", "EXJPYmHNbNAYmZvMNkHEhrBrJ", "http://www.QmQKaaZ.info/", "Schimmel Company", false, 744290621, 74771059},
	{"8c514a6d-3363-453e-a427-023b7fae1142", "MarWsYNWANqXlYuHFWIuFoNvX", "tccXlnYclQIjjbaZbaOBwgtVt", "http://www.OxDkypZ.net/", "Beahan Agency", false, 1206575362, 159350386},
	{"a8c976cb-f6bb-46f5-aca9-1d78d33c7325", "qnRpvcbkOoAsHdrvcVvsNTvbX", "kqgRDpcJYCigLbqSQFGFPHfMG", "http://BPXcHVW.biz/", "Will Ltd", false, 243980551, 616638567},
}

var events = []InternalEvent{
	{"98ca8b18-0289-4125-a88c-1a59097acb21", "Perferendis sit accusantium aut voluptatem consequatur.", model.ProtocolTypeConnection, model.EventTypeAction, "8c514a6d-3363-453e-a427-023b7fae1142", 13721960},
	{"0130b796-7afb-435a-9235-d6b03ad9273b", "Accusantium consequatur voluptatem aut sit perferendis.", model.ProtocolTypeNone, model.EventTypeNotification, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 38274631},
	{"21dbd8a3-d13e-4918-a8c9-ef35ad35fe06", "Aut accusantium sit perferendis consequatur voluptatem.", model.ProtocolTypeConnection, model.EventTypeNotification, "8c514a6d-3363-453e-a427-023b7fae1142", 80228851},
	{"284b2029-f8b6-4e0d-8461-20f3ccd6da76", "Accusantium voluptatem perferendis consequatur aut sit.", model.ProtocolTypeConnection, model.EventTypeAction, "bd5b6b66-a2cd-451f-b906-37ea3dbb1301", 124197655},
	{"156af804-a91e-4b12-9eba-d5ca31fee257", "Voluptatem aut perferendis sit accusantium consequatur.", model.ProtocolTypeNone, model.EventTypeAction, "8c514a6d-3363-453e-a427-023b7fae1142", 153865357},
	{"94b7d52b-8aaf-4cb6-b829-daea4c7b072e", "Sit voluptatem aut consequatur accusantium perferendis.", model.ProtocolTypeConnection, model.EventTypeAction, "95b2eab7-d664-4868-b684-3dbbcc3e0375", 172735825},
	{"af2b1eff-b194-48d0-a28d-45581976d990", "Sit consequatur accusantium aut perferendis voluptatem.", model.ProtocolTypeConnection, model.EventTypeAction, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 196384182},
	{"2eac5b38-5b96-4f79-80d2-0bc4930f9a78", "Sit perferendis accusantium voluptatem consequatur aut.", model.ProtocolTypeProof, model.EventTypeNotification, "95b2eab7-d664-4868-b684-3dbbcc3e0375", 227207328},
	{"40846661-d9e8-44c1-9eab-75a38be44b0b", "Accusantium sit voluptatem perferendis consequatur aut.", model.ProtocolTypeProof, model.EventTypeNotification, "bd5b6b66-a2cd-451f-b906-37ea3dbb1301", 247518927},
	{"8b9c26ad-fe27-4c37-854b-99d05777c556", "Consequatur sit accusantium perferendis voluptatem aut.", model.ProtocolTypeNone, model.EventTypeAction, "95b2eab7-d664-4868-b684-3dbbcc3e0375", 254393056},
	{"deb0aed0-24f7-4ba1-b39c-1e1713af8bb0", "Consequatur accusantium aut perferendis voluptatem sit.", model.ProtocolTypeProof, model.EventTypeNotification, "bd5b6b66-a2cd-451f-b906-37ea3dbb1301", 270519173},
	{"a92f8949-e783-4e40-b2b7-20fd210800c1", "Perferendis voluptatem sit aut consequatur accusantium.", model.ProtocolTypeNone, model.EventTypeNotification, "8c514a6d-3363-453e-a427-023b7fae1142", 319418545},
	{"ef7ee797-ab87-4f76-aff3-b40c43fe7b2a", "Accusantium aut voluptatem perferendis sit consequatur.", model.ProtocolTypeNone, model.EventTypeNotification, "95b2eab7-d664-4868-b684-3dbbcc3e0375", 328155582},
	{"c4fad4c6-924f-44eb-bd6a-5445eda85d5c", "Sit perferendis voluptatem accusantium aut consequatur.", model.ProtocolTypeCredential, model.EventTypeAction, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 380591732},
	{"59d5ee28-e53a-41ee-abf1-8921ebc45b70", "Perferendis voluptatem sit accusantium consequatur aut.", model.ProtocolTypeProof, model.EventTypeNotification, "95b2eab7-d664-4868-b684-3dbbcc3e0375", 390548840},
	{"f90b35f9-bf4e-4ff1-b2f7-d7d0011b22ba", "Voluptatem accusantium consequatur perferendis sit aut.", model.ProtocolTypeProof, model.EventTypeAction, "95b2eab7-d664-4868-b684-3dbbcc3e0375", 424125698},
	{"de7bb40e-086d-442e-a267-6e2ab0107a3d", "Aut voluptatem sit consequatur perferendis accusantium.", model.ProtocolTypeProof, model.EventTypeAction, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 424815312},
	{"223ea1dc-5efb-4cb1-8bca-ea22443ecaed", "Accusantium perferendis consequatur voluptatem sit aut.", model.ProtocolTypeConnection, model.EventTypeNotification, "bd5b6b66-a2cd-451f-b906-37ea3dbb1301", 434127563},
	{"64c05824-f387-43e8-ae56-97605577a0e1", "Consequatur aut perferendis sit voluptatem accusantium.", model.ProtocolTypeConnection, model.EventTypeAction, "95b2eab7-d664-4868-b684-3dbbcc3e0375", 438310759},
	{"076dca94-0753-43af-b479-24f064dd1716", "Perferendis consequatur accusantium sit aut voluptatem.", model.ProtocolTypeProof, model.EventTypeAction, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 553285558},
	{"5536ed2f-1eed-4114-92a6-14fc44f131ca", "Accusantium perferendis aut voluptatem consequatur sit.", model.ProtocolTypeCredential, model.EventTypeNotification, "8c514a6d-3363-453e-a427-023b7fae1142", 597059561},
	{"882fdc12-84aa-4336-90bd-d0c5c2e46881", "Aut accusantium consequatur voluptatem sit perferendis.", model.ProtocolTypeCredential, model.EventTypeNotification, "95b2eab7-d664-4868-b684-3dbbcc3e0375", 652863386},
	{"8ecf07d3-ae75-4043-8306-9ccfc9deb388", "Consequatur aut voluptatem accusantium sit perferendis.", model.ProtocolTypeBasicMessage, model.EventTypeNotification, "8c514a6d-3363-453e-a427-023b7fae1142", 666905482},
	{"1bab6184-75a8-4c69-bd4d-0ddc680e1279", "Aut sit voluptatem consequatur accusantium perferendis.", model.ProtocolTypeBasicMessage, model.EventTypeAction, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 678941190},
	{"ac12d8cc-42f0-4b42-aa0b-35ae1cf8be22", "Aut sit consequatur perferendis accusantium voluptatem.", model.ProtocolTypeBasicMessage, model.EventTypeAction, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 702635374},
	{"753037cc-1daf-47ba-bf7a-b9ca209c6a66", "Consequatur aut accusantium sit perferendis voluptatem.", model.ProtocolTypeConnection, model.EventTypeNotification, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 759855619},
	{"c96ffa1b-e8f9-4a79-afdc-ff4b1bdceb33", "Voluptatem consequatur aut perferendis accusantium sit.", model.ProtocolTypeProof, model.EventTypeAction, "8c514a6d-3363-453e-a427-023b7fae1142", 769802518},
	{"acd5902f-90b1-4ad7-b309-e2242ee1889e", "Consequatur perferendis sit voluptatem accusantium aut.", model.ProtocolTypeNone, model.EventTypeAction, "8c514a6d-3363-453e-a427-023b7fae1142", 786277538},
	{"3a70a829-d337-4ac8-97e7-3e43c2412a17", "Voluptatem consequatur perferendis accusantium aut sit.", model.ProtocolTypeConnection, model.EventTypeAction, "8c514a6d-3363-453e-a427-023b7fae1142", 847327160},
	{"eb9475be-d6c8-464f-906b-f99c82fea4c6", "Sit accusantium consequatur perferendis aut voluptatem.", model.ProtocolTypeConnection, model.EventTypeNotification, "bd5b6b66-a2cd-451f-b906-37ea3dbb1301", 911002377},
	{"75920851-99e6-4c42-ab46-6c7ea48a3e69", "Consequatur voluptatem aut accusantium sit perferendis.", model.ProtocolTypeBasicMessage, model.EventTypeNotification, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 956827148},
	{"4da589f5-2a55-4168-bc36-91d08725aa77", "Perferendis aut sit consequatur accusantium voluptatem.", model.ProtocolTypeProof, model.EventTypeNotification, "bd5b6b66-a2cd-451f-b906-37ea3dbb1301", 991987799},
	{"58718267-84a1-46d4-b8d9-0e32a14dc9b2", "Aut consequatur accusantium sit voluptatem perferendis.", model.ProtocolTypeProof, model.EventTypeAction, "95b2eab7-d664-4868-b684-3dbbcc3e0375", 994460583},
	{"6ec55145-f4ac-431b-b253-caa7c55f3b02", "Accusantium aut voluptatem consequatur sit perferendis.", model.ProtocolTypeNone, model.EventTypeAction, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 999543724},
	{"96deb528-6e81-4e10-9a7e-93782bfcd67c", "Sit consequatur voluptatem accusantium perferendis aut.", model.ProtocolTypeConnection, model.EventTypeAction, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 1044784344},
	{"924fc58b-8013-4e92-9c7b-a6ca663a75f9", "Consequatur aut sit voluptatem accusantium perferendis.", model.ProtocolTypeProof, model.EventTypeNotification, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 1076232012},
	{"a24751a5-20f0-4dc0-a5ec-4d3486051e26", "Accusantium sit perferendis consequatur aut voluptatem.", model.ProtocolTypeCredential, model.EventTypeNotification, "bd5b6b66-a2cd-451f-b906-37ea3dbb1301", 1098615336},
	{"3431c405-b277-48bd-a802-92adf31b4fc9", "Voluptatem accusantium consequatur sit aut perferendis.", model.ProtocolTypeConnection, model.EventTypeNotification, "95b2eab7-d664-4868-b684-3dbbcc3e0375", 1136072426},
	{"cb3954de-43bb-4a41-9c6e-d5cf292bad32", "Aut consequatur voluptatem accusantium sit perferendis.", model.ProtocolTypeCredential, model.EventTypeAction, "95b2eab7-d664-4868-b684-3dbbcc3e0375", 1166733839},
	{"095fd5dd-720d-4177-8741-f654686eab36", "Voluptatem accusantium aut perferendis consequatur sit.", model.ProtocolTypeCredential, model.EventTypeAction, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 1217033997},
	{"ac4bcea4-a19d-4a3d-ac16-e61f4b80430d", "Voluptatem consequatur perferendis accusantium aut sit.", model.ProtocolTypeProof, model.EventTypeNotification, "8c514a6d-3363-453e-a427-023b7fae1142", 1224592555},
	{"f2c9e355-95c7-43f3-89f7-739c0dc24005", "Sit accusantium voluptatem consequatur aut perferendis.", model.ProtocolTypeConnection, model.EventTypeAction, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 1264006456},
	{"5c22bf69-28b0-4ed2-9693-0bc19a9cf2c7", "Voluptatem consequatur perferendis sit accusantium aut.", model.ProtocolTypeProof, model.EventTypeAction, "8c514a6d-3363-453e-a427-023b7fae1142", 1272905241},
	{"6a760732-39a5-4534-a546-09966ede0358", "Accusantium consequatur aut perferendis voluptatem sit.", model.ProtocolTypeConnection, model.EventTypeAction, "bd5b6b66-a2cd-451f-b906-37ea3dbb1301", 1325639313},
	{"d164685e-7a39-49f7-a40c-ad0d710242e0", "Voluptatem accusantium aut perferendis sit consequatur.", model.ProtocolTypeNone, model.EventTypeAction, "95b2eab7-d664-4868-b684-3dbbcc3e0375", 1346400138},
	{"dbda2e65-9d3c-4777-8806-b9b1453ee7f4", "Sit perferendis consequatur aut voluptatem accusantium.", model.ProtocolTypeProof, model.EventTypeNotification, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 1346979259},
	{"d243006e-59bf-4a9c-b842-93088289c974", "Aut consequatur voluptatem sit perferendis accusantium.", model.ProtocolTypeProof, model.EventTypeNotification, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 1473588752},
	{"45cc0499-99ab-4349-9c41-ce7ea2e744bc", "Sit voluptatem perferendis aut accusantium consequatur.", model.ProtocolTypeConnection, model.EventTypeNotification, "a8c976cb-f6bb-46f5-aca9-1d78d33c7325", 1497627274},
	{"522fcc28-e89a-4971-9d3c-97fb1dcac492", "Perferendis accusantium consequatur aut voluptatem sit.", model.ProtocolTypeProof, model.EventTypeNotification, "8c514a6d-3363-453e-a427-023b7fae1142", 1524954993},
	{"d6222ecc-5d9f-49b2-97b3-ca0c7b656a8d", "Perferendis sit voluptatem consequatur aut accusantium.", model.ProtocolTypeConnection, model.EventTypeAction, "bd5b6b66-a2cd-451f-b906-37ea3dbb1301", 1575581128},
}
