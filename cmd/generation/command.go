package generation

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	dbType        string
	host          string
	user          string
	password      string
	schema        string
	dbname        string
	serviceName   string
	packageName   string
	goPackageName string
	ignoreTables  []string
	ignoreColumns []string
	fieldStyle    string
	table         string
	port          int
)

var GenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generates protobuf",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := generateSchema(table, ignoreTables, ignoreColumns, serviceName, fieldStyle, dbType)
		if nil != err {
			log.Fatal(err)
		}

		if nil != s {
			fmt.Println(s)
		}
	},
}

func init() {
	GenCmd.Flags().StringVarP(&dbType, "db_type", "", "mysql", "the database type. mysql | postgres")
	GenCmd.Flags().StringVarP(&host, "host", "", "localhost", "the database host")
	GenCmd.Flags().IntVarP(&port, "port", "", 3306, "the database port")
	GenCmd.Flags().StringVarP(&user, "user", "", "root", "the database user")
	GenCmd.Flags().StringVarP(&password, "password", "", "", "the database password")
	GenCmd.Flags().StringVarP(&schema, "schema", "", "", "the database schema")
	GenCmd.Flags().StringVarP(&dbname, "dbname", "", "", "the database name")
	GenCmd.Flags().StringVarP(&table, "table", "", "", "the table schema. multiple tables ',' split. ")
	GenCmd.Flags().StringVarP(&serviceName, "service_name", "", schema, "the protocol buffer package. defaults to the database schema.")
	GenCmd.Flags().StringVarP(&packageName, "package", "", schema, "the protocol buffer package. defaults to the database schema.")
	GenCmd.Flags().StringVarP(&goPackageName, "go_package", "", "", "the protocol buffer go_package. defaults to the database schema.")
	GenCmd.Flags().StringSliceVarP(&ignoreTables, "ignore_tables", "", []string{}, "a comma spaced list of tables to ignore")
	GenCmd.Flags().StringSliceVarP(&ignoreColumns, "ignore_columns", "", []string{}, "a comma spaced list of mysql columns to ignore")
	GenCmd.Flags().StringVarP(&fieldStyle, "field_style", "", "sql_pb", "gen protobuf field style. sql_pb | sqlPb")

}
