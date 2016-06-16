
# Test Data 1

How to create soft-link

```bash
echo "Hello" >tmpfile
ln -s tmpfile linkedfile
ls -l linkedfile | grep '^l'
rm linkedfile tmpfile
```

Remove file

```bash
echo "Hello again" > tmpfile2
rm tmpfile2
```
