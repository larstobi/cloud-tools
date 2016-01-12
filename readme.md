# Cloud tools

Project containing several command line utilities that are useful when using 
Apache Cloudstack and/or Hashicorp Terraform:


## Tools:

### cloudstack-templates
 
     cloudstack-templates [<keyword>]
 
 Will list cloudstack templates sorted by date, and optionally filtered by keyword
  
### cloudstack-vpn 
  
    cloudstack-vpn <vpc name>
    
Will enable vpn for given vpc-name

### terraform-wrapper

    terraform-wrapper [args.....]
    
Will call _terraform_ passing environment variables found in cloud-config.yml
 
 
## Configuration

Utilties will look for a file called _cloud-config.yml_ in cwd containing references to passwords stored in 
a pass password store and inline clear text variables
