package sofar

import (
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/ports"
)

type field struct {
	register  int
	name      string
	valueType string
	factor    string
	unit      string
}

type registerRange struct {
	start       int
	end         int
	replyFields []field
}

var allRegisterRanges = []registerRange{
	rrGridOutput,
	rrPVOutput,
	rrEnergyTodayTotals,
	rrSystemInfo,
	rrBatOutput,
	rrRatio,
}

func GetAllRegisterNames() []string {
	result := make([]string, 0)
	for _, rr := range allRegisterRanges {
		for _, f := range rr.replyFields {
			if f.name == "" || f.valueType == "" {
				// Measurements without a name or value type are ignored in replies
				continue
			}
			result = append(result, f.name)
		}
	}
	return result
}

func getDiscoveryFields(nameFilter func(string) bool) []ports.DiscoveryField {
	result := make([]ports.DiscoveryField, 0)
	for _, rr := range allRegisterRanges {
		for _, f := range rr.replyFields {
			if f.name != "" && f.valueType != "" && nameFilter(f.name) {
				result = append(result, ports.DiscoveryField{Name: f.name, Factor: f.factor, Unit: f.unit})
			}
		}
	}
	return result
}

var rrSystemInfo = registerRange{
	start: 0x400,
	end:   0x43a,
	replyFields: []field{
		{0x0404, "SysState", "U16", "", ""},
		{0x0405, "Fault1", "U16", "", ""},
		{0x0406, "Fault2", "U16", "", ""},
		{0x0407, "Fault3", "U16", "", ""},
		{0x0408, "Fault4", "U16", "", ""},
		{0x0409, "Fault5", "U16", "", ""},
		{0x040A, "Fault6", "U16", "", ""},
		{0x040B, "Fault7", "U16", "", ""},
		{0x040C, "Fault8", "U16", "", ""},
		{0x040D, "Fault9", "U16", "", ""},
		{0x040E, "Fault10", "U16", "", ""},
		{0x040F, "Fault11", "U16", "", ""},
		{0x0410, "Fault12", "U16", "", ""},
		{0x0411, "Fault13", "U16", "", ""},
		{0x0412, "Fault14", "U16", "", ""},
		{0x0413, "Fault15", "U16", "", ""},
		{0x0414, "Fault16", "U16", "", ""},
		{0x0415, "Fault17", "U16", "", ""},
		{0x0416, "Fault18", "U16", "", ""},
		{0x0417, "Countdown", "U16", "1", "seconds"},
		{0x0418, "Temperature_Env1", "I16", "1", "℃"},
		{0x0419, "Temperature_Env2", "I16", "1", "℃"},
		{0x041A, "Temperature_HeatSink1", "I16", "1", "℃"},
		{0x041B, "Temperature_HeatSink2", "I16", "1", "℃"},
		{0x041C, "Temperature_HeatSink3", "I16", "1", "℃"},
		{0x041D, "Temperature_HeatSink4", "I16", "1", "℃"},
		{0x041E, "Temperature_HeatSink5", "I16", "1", "℃"},
		{0x041F, "Temperature_HeatSink6", "I16", "1", "℃"},
		{0x0420, "Temperature_Inv1", "I16", "1", "℃"},
		{0x0421, "Temperature_Inv2", "I16", "1", "℃"},
		{0x0422, "Temperature_Inv3", "I16", "1", "℃"},
		{0x0423, "Temp_Rsvd1", "I16", "1", "℃"},
		{0x0424, "Temp_Rsvd2", "I16", "1", "℃"},
		{0x0425, "Temp_Rsvd3", "I16", "1", "℃"},
		{0x0426, "GenerationTime_Today", "U16", "1", "min"}, // HA uses d, h, min, s not Minute
		{0x0427, "GenerationTime_Total", "U32", "1", "min"},
		{0x0428, "", "", "", ""},
		{0x0429, "ServiceTime_Total", "U32", "1", "min"},
		{0x042A, "", "", "", ""},
		{0x042B, "InsulationResistance", "U16", "1", "kΩ"},
		{0x042C, "SysTime_Year", "U16", "", ""},
		{0x042D, "SysTime_Month", "U16", "", ""},
		{0x042E, "SysTime_Date", "U16", "1", "d"},
		{0x042F, "SysTime_Hour", "U16", "1", "h"},
		{0x0430, "SysTime_Minute", "U16", "1", "min"},
		{0x0431, "SysTime_Second", "U16", "1", "s"},
		{0x0432, "Fault19", "U16", "", ""},
		{0x0433, "Fault20", "U16", "", ""},
		{0x0434, "Fault21", "U16", "", ""},
		{0x0435, "Fault22", "U16", "", ""},
		{0x0436, "Fault23", "U16", "", ""},
		{0x0437, "Fault24", "U16", "", ""},
		{0x0438, "Fault25", "U16", "", ""},
		{0x0439, "Fault26", "U16", "", ""},
		{0x043A, "Fault27", "U16", "", ""},
	},
}
var rrEnergyTodayTotals = registerRange{
	start: 0x680,
	end:   0x69B,
	replyFields: []field{
		{0x684, "PV_Generation_Today", "U32", "0.01", "kWh"},
		{0x686, "PV_Generation_Total", "U32", "0.1", "kWh"},
		{0x688, "Load_Consumption_Today", "U32", "0.01", "kWh"},
		{0x68A, "Load_Consumption_Total", "U32", "0.1", "kWh"},
		{0x68C, "Energy_Purchase_Today", "U32", "0.01", "kWh"},
		{0x68E, "Energy_Purchase_Total", "U32", "0.1", "kWh"},
		{0x690, "Energy_Selling_Today", "U32", "0.01", "kWh"},
		{0x692, "Energy_Selling_Total", "U32", "0.1", "kWh"},
		{0x694, "Bat_Charge_Today", "U32", "0.01", "kWh"},
		{0x696, "Bat_Charge_Total", "U32", "0.1", "kWh"},
		{0x698, "Bat_Discharge_Today", "U32", "0.01", "kWh"},
		{0x69A, "Bat_Discharge_Total", "U32", "0.1", "kWh"},
	},
}

var rrPVOutput = registerRange{
	start: 0x580,
	end:   0x589,
	replyFields: []field{
		{0x0584, "Voltage_PV1", "U16", "0.1", "V"},
		{0x0585, "Current_PV1", "U16", "0.01", "A"},
		{0x0586, "Power_PV1", "U16", "0.01", "kW"},
		{0x0587, "Voltage_PV2", "U16", "0.1", "V"},
		{0x0588, "Current_PV2", "U16", "0.01", "A"},
		{0x0589, "Power_PV2", "U16", "0.01", "kW"},
	},
}
var rrGridOutput = registerRange{
	start: 0x480,
	end:   0x4bc,
	replyFields: []field{
		{0x0484, "Frequency_Grid", "U16", "0.01", "Hz"},
		{0x0485, "ActivePower_Output_Total", "I16", "0.01", "kW"},
		{0x0486, "ReactivePower_Output_Total", "I16", "0.01", "kW"},
		{0x0487, "ApparentPower_Output_Total", "I16", "0.01", "kW"},
		{0x0488, "ActivePower_PCC_Total", "I16", "0.01", "kW"},
		{0x0489, "ReactivePower_PCC_Total", "I16", "0.01", "kW"},
		{0x048A, "ApparentPower_PCC_Total", "I16", "0.01", "kW"},
		{0x048B, "GridOutput_Rsvd1", "", "", ""},
		{0x048C, "GridOutput_Rsvd2", "", "", ""},
		{0x048D, "Voltage_Phase_R", "U16", "0.1", "V"},
		{0x048E, "Current_Output_R", "U16", "0.01", ""},
		{0x048F, "ActivePower_Output_R", "I16", "0.01", "kW"},
		{0x0490, "ReactivePower_Output_R", "I16", "0.01", "kW"},
		{0x0491, "PowerFactor_Output_R", "I16", "0.001", "p.u."},
		{0x0492, "Current_PCC_R", "U16", "0.01", ""},
		{0x0493, "ActivePower_PCC_R", "I16", "0.01", "kW"},
		{0x0494, "ReactivePower_PCC_R", "I16", "0.01", "kW"},
		{0x0495, "PowerFactor_PCC_R", "I16", "0.001", "p.u."},
		{0x0496, "R_Rsvd1", "", "", ""},
		{0x0497, "R_Rsvd2", "", "", ""},
		{0x0498, "Voltage_Phase_S", "U16", "0.1", "V"},
		{0x0499, "Current_Output_S", "U16", "0.01", ""},
		{0x049A, "ActivePower_Output_S", "I16", "0.01", "kW"},
		{0x049B, "ReactivePower_Output_S", "I16", "0.01", "kW"},
		{0x049C, "PowerFactor_Output_S", "I16", "0.001", "p.u."},
		{0x049D, "Current_PCC_S", "U16", "0.01", ""},
		{0x049E, "ActivePower_PCC_S", "I16", "0.01", "kW"},
		{0x049F, "ReactivePower_PCC_S", "I16", "0.01", "kW"},
		{0x04A0, "PowerFactor_PCC_S", "I16", "0.001", "p.u."},
		{0x04A1, "S_Rsvd1", "", "", ""},
		{0x04A2, "S_Rsvd2", "", "", ""},
		{0x04A3, "Voltage_Phase_T", "U16", "0.1", "V"},
		{0x04A4, "Current_Output_T", "U16", "0.01", ""},
		{0x04A5, "ActivePower_Output_T", "I16", "0.01", "kW"},
		{0x04A6, "ReactivePower_Output_T", "I16", "0.01", "kW"},
		{0x04A7, "PowerFactor_Output_T", "I16", "0.001", "p.u."},
		{0x04A8, "Current_PCC_T", "U16", "0.01", ""},
		{0x04A9, "ActivePower_PCC_T", "I16", "0.01", "kW"},
		{0x04AA, "ReactivePower_PCC_T", "I16", "0.01", "kW"},
		{0x04AB, "PowerFactor_PCC_T", "I16", "0.001", "p.u."},
		{0x04AC, "T_Rsvd1", "", "", ""},
		{0x04AD, "T_Rsvd2", "", "", ""},
		{0x04AE, "ActivePower_PV_Ext", "U16", "0.01", "kW"},
		{0x04AF, "ActivePower_Load_Sys", "U16", "0.01", "kW"},
		{0x04B0, "Voltage_Phase_L1N", "U16", "0.1", "V"},
		{0x04B1, "Current_Output_L1N", "U16", "0.01", ""},
		{0x04B2, "ActivePower_Output_L1N", "I16", "0.01", "kW"},
		{0x04B3, "Current_PCC_L1N", "U16", "0.01", ""},
		{0x04B4, "ActivePower_PCC_L1N", "I16", "0.01", "kW"},
		{0x04B5, "Voltage_Phase_L2N", "U16", "0.1", "V"},
		{0x04B6, "Current_Output_L2N", "U16", "0.01", ""},
		{0x04B7, "ActivePower_Output_L2N", "I16", "0.01", "kW"},
		{0x04B8, "Current_PCC_L2N", "U16", "0.01", ""},
		{0x04B9, "ActivePower_PCC_L2N", "I16", "0.01", "kW"},
		{0x04BA, "Voltage_Line_L1", "U16", "0.1", "V"},
		{0x04BB, "Voltage_Line_L2", "U16", "0.1", "V"},
		{0x04BC, "Voltage_Line_L3", "U16", "0.1", "V"},
	},
}

var rrBatOutput = registerRange{
	start: 0x600,
	end:   0x611,
	replyFields: []field{
		{0x0604, "Voltage_Bat1", "U16", "0.1", "V"},
		{0x0605, "Current_Bat1", "I16", "0.01", "A"},
		{0x0606, "Power_Bat1", "I16", "0.01", "kW"},
		{0x0607, "Temperature_Env_Bat1", "I16", "1", "℃"},
		{0x0608, "SOC_Bat1", "U16", "1", "%"},
		{0x0609, "SOH_Bat1", "U16", "1", "%"},
		{0x060A, "ChargeCycle_Bat1", "U16", "1", ""},
		{0x060B, "Voltage_Bat2", "U16", "0.1", "V"},
		{0x060C, "Current_Bat2", "I16", "0.01", "A"},
		{0x060D, "Power_Bat2", "I16", "0.01", "kW"},
		{0x060E, "Temperature_Env_Bat2", "I16", "1", "℃"},
		{0x060F, "SOC_Bat2", "U16", "1", "%"},
		{0x0610, "SOH_Bat2", "U16", "1", "%"},
		{0x0611, "ChargeCycle_Bat2", "U16", "1", ""},
	},
}

var rrRatio = registerRange{
	start: 0x1030,
	end:   0x103D,
	replyFields: []field{
		{0x1039, "PV_Generation_Ratio", "U16", "0.001", ""},
		{0x103A, "Energy_Purchase_Ratio", "U16", "0.001", ""},
		{0x103B, "Energy_Selling_Ratio", "U16", "0.001", ""},
		{0x103C, "Bat_Charge_Ratio", "U16", "0.001", ""},
		{0x103D, "Bat_Discharge_Ratio", "U16", "0.001", ""},
	},
}
