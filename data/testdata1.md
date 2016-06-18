
# Test Data 1

How to create soft-link

```bashtest
echo "Hello" >tmpfile
ln -s tmpfile linkedfile
ls -l linkedfile | grep '^l'
rm linkedfile tmpfile
```

Failed test

```bashtest
rm tmpfile
```

No test runner

```bash
ls .
```
