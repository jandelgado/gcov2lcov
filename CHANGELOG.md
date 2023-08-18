# changelog for gcov2lcov

## 1.0.6 [2023-08-18]

* performance otimizations (thanks to zzh8829, #16)
* dependency upgrades

## 1.0.5 [2021-04-28]

* new: option  `-use-absolute-source-path` - when set absolut path names are
       used for the SF value (#10)
* compile with go 1.16

## 1.0.2 [2020-04-25]

* fix calculation of LH and LF values which led to wrong calculation of
  test coverage in coveralls

## 1.0.1 [2020-04-25]

* avoid duplicate DA records for same lines (see
  https://github.com/jandelgado/gcov2lcov-action/issues/2)

## 1.0.0 [2019-10-07]

* initial release
