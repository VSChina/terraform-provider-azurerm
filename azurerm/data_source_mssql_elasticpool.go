package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMsSqlElasticpool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMsSqlElasticpoolRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"max_size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_size_gb": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"per_db_min_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"per_db_max_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),

			"zone_redundant": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceArmMsSqlElasticpoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mssql.ElasticPoolsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	elasticPoolName := d.Get("name").(string)
	serverName := d.Get("server_name").(string)

	resp, err := client.Get(ctx, resourceGroup, serverName, elasticPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Elasticpool %q (Resource Group %q, SQL Server %q) was not found", elasticPoolName, resourceGroup, serverName)
		}

		return fmt.Errorf("Error making Read request on AzureRM Elasticpool %s (Resource Group %q, SQL Server %q): %+v", elasticPoolName, resourceGroup, serverName, err)
	}

	if id := resp.ID; id != nil {
		d.SetId(*resp.ID)
	}
	d.Set("name", elasticPoolName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("server_name", serverName)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	flattenAndSetTags(d, resp.Tags)

	if props := resp.ElasticPoolProperties; props != nil {
		d.Set("max_size_gb", float64(*props.MaxSizeBytes/int64(1073741824)))
		d.Set("max_size_bytes", props.MaxSizeBytes)

		d.Set("zone_redundant", props.ZoneRedundant)

		if perDbSettings := props.PerDatabaseSettings; perDbSettings != nil {
			d.Set("per_db_min_capacity", perDbSettings.MinCapacity)
			d.Set("per_db_max_capacity", perDbSettings.MaxCapacity)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}
