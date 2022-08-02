# optify 

Turning yaml/json file to command line long opts.

## usage 

Suppose you have a `cmd` which takes A LOT of command line options. 
But you only change a few of them each time you invoke `cmd`.
So you want to store the default options in a yaml file `defaults.yaml`
and only specify the overriding opts in command line.

You can do the following with `optify`:

`optify defaults.yaml -- cmd subcmd1 subcmd2 --opt1 o1 --opt2 o2`

To much to type? alias it

```bash
alias cmd="optify defaults.yaml -- cmd"
cmd subcmd1 --opt1 o1 --opt2 o2
```

e.g. to manage your GCE instance definition as yaml file

```bash
optify vm_defaults.yaml -- gcloud compute instance create
```
