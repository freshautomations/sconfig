package cmd

import (
	"fmt"
	"github.com/freshautomations/sconfig/defaults"
	"github.com/freshautomations/sconfig/exit"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
)

const epsilon = 1e-9 // Margin of error

type FlagsType struct {
	Strict bool
	Type   string
}

var inputFlags FlagsType

func CheckArgs(cmd *cobra.Command, args []string) error {
	validateArgs := cobra.MinimumNArgs(2)
	if err := validateArgs(cmd, args); err != nil {
		return err
	}

	for i := 1; i < len(args); i++ {
		if (len(strings.Split(args[i], "="))) < 2 {
			return fmt.Errorf("key=value not present in %s", args[i])
		}

	}

	fileName := args[0]
	_, err := os.Stat(fileName)
	return err
}

func translateUintSlice(incoming []string) (result []uint64, err error) {
	for _, v := range incoming {
		translated, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, translated)
	}
	return result, nil
}

func RunRoot(cmd *cobra.Command, args []string) (output string, err error) {
	fileName := args[0]

	viper.SetConfigFile(fileName)
	err = viper.ReadInConfig()
	if err != nil {
		if _, IsUnsupportedExtension := err.(viper.UnsupportedConfigError); IsUnsupportedExtension {
			viper.SetConfigType("toml")
			err = viper.ReadInConfig()
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	for i := 1; i < len(args); i++ {
		keyvalue := strings.Split(args[i], "=")
		key := keyvalue[0]
		value := keyvalue[1]
		keyvaluetype := inputFlags.Type

		if inputFlags.Strict && !viper.IsSet(key) {
			return "", fmt.Errorf("key does not exist in %s", args[i])
		}

		if keyvaluetype == "" {
			switch viper.Get(key).(type) {
			case bool:
				keyvaluetype = "bool"
			case float32, float64:
				keyvaluetype = "float"
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
				keyvaluetype = "int"
			case []int:
				keyvaluetype = "intSlice"
			case []string:
				keyvaluetype = "stringSlice"
			case time.Time:
				keyvaluetype = "Time"
			case time.Duration:
				keyvaluetype = "Duration"
			case string:
				keyvaluetype = "string"
			// section
			case map[string]interface{}, map[string]string:
				return "", fmt.Errorf("cannot set section or map types in %s", args[i])
			default:
				keyvaluetype = "string"
			}
		}

		switch strings.ToLower(keyvaluetype) {
		case "bool", "boolean":
			boolvalue, err := strconv.ParseBool(value)
			if err != nil {
				return "", err
			}
			viper.Set(key, boolvalue)
		case "float", "float32", "float64":
			floatvalue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return "", err
			}
			viper.Set(key, floatvalue)
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
			uintvalue, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				intvalue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return "", err
				} else {
					viper.Set(key, intvalue)
				}
			} else {
				viper.Set(key, uintvalue)
			}
		case "intslice":
			s := strings.TrimPrefix(value, "[")
			s = strings.TrimSuffix(s, "]")
			result1 := strings.Split(s, ",")
			result2 := strings.Split(s, " ")
			if len(result2) > len(result1) {
				translated, err := translateUintSlice(result2)
				if err != nil {
					return "", err
				}
				viper.Set(key, translated)
			} else {
				translated, err := translateUintSlice(result1)
				if err != nil {
					return "", err
				}
				viper.Set(key, translated)
			}
		case "stringslice":
			s := strings.TrimPrefix(value, "[")
			s = strings.TrimSuffix(s, "]")
			result1 := strings.Split(s, ",")
			result2 := strings.Split(s, " ")
			if len(result2) > len(result1) {
				viper.Set(key, result2)
			} else {
				viper.Set(key, result1)
			}
		case "time":
			result, err := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", value)
			if err != nil {
				return "", err
			}
			viper.Set(key, result)
		case "duration":
			result, err := time.ParseDuration(value)
			if err != nil {
				return "", err
			}
			viper.Set(key, result)
		case "string", "str":
			viper.Set(key, value)
		default:
			return "", fmt.Errorf("invalid type in %s", args[i])
		}
	}
	if len(args) == 2 {
		return "1 configuration item set\n", viper.WriteConfig()
	}
	return fmt.Sprintf("%d configuration items set\n", len(args)-1), viper.WriteConfig()

}

func runRootWrapper(cmd *cobra.Command, args []string) {
	if result, err := RunRoot(cmd, args); err != nil {
		exit.Fail(err)
	} else {
		exit.Succeed(result)
	}
}

func Execute() error {
	var rootCmd = &cobra.Command{
		Version: defaults.Version,
		Use:     "sconfig",
		Short:   "SCONFIG - simple configuration changer for Shell",
		Long: `Configuration file changer for the Linux Shell.
Source and documentation is available at https://github.com/freshautomations/sconfig`,
		Args: CheckArgs,
		Run:  runRootWrapper,
	}
	rootCmd.Use = "sconfig <filename> <key=value> [<key=value>] ..."

	pflag.BoolVarP(&inputFlags.Strict, "strict", "s", false, "Only allow changes but not new entries.")
	pflag.StringVarP(&inputFlags.Type, "type", "t", "", "Override value(s) type. All values will have the same type. (int, float, bool, string, stringSlice, intSlice) Slices are comma- or space-separated. Optionally, slices can be in brackets: [1,2,3]")

	return rootCmd.Execute()
}
