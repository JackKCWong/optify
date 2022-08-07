# optify 

Turning yaml/json file to command line long opts.

## usage 

Suppose you have a `cmd` which takes A LOT of command line options. 
But you only change a few of them each time you invoke `cmd`.
So you want to store the default options in a yaml file `defaults.yaml`
and only specify the overriding opts in command line.

You can do the following with `optify`:

`optify cmd --opt1 o1 --opt2 o2 -- defaults.yaml`

You can alias it:

```bash
alias cmd="optify cmd"
cmd --opt1 o1 --opt2 o2 -- defaults.yaml
```


e.g. to manage your GCE instance definition as yaml file

```bash
optify gcloud compute instance create -- vm_defaults.yaml 
# or 
alias gcloud="optify gcloud"
gcloud compute instance create --name foo -- defaults.yaml
gcloud compute instance create --name foo # without '-- defaults.yaml' it just works like invoking the command directly
```
