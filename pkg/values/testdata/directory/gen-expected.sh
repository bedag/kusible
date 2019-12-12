if ! [ -d expected ]; then
  echo "Script must be execuded in the testdata/directory directory!"
  exit 1
fi

TESTCASE=multi-dir-multi-file
PERMUTATION="123"
spruce merge \
  "${TESTCASE}"/multi-file-01/file-01.yml \
  "${TESTCASE}"/multi-file-01/file-02.yml \
  "${TESTCASE}"/multi-file-01/file-03.yml \
  "${TESTCASE}"/multi-file-02/file-01.yml \
  "${TESTCASE}"/multi-file-02/file-02.yml \
  "${TESTCASE}"/multi-file-02/file-03.yml \
  "${TESTCASE}"/multi-file-03/file-01.yml \
  "${TESTCASE}"/multi-file-03/file-02.yml \
  "${TESTCASE}"/multi-file-03/file-03.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="213"
spruce merge \
  "${TESTCASE}"/multi-file-02/file-01.yml \
  "${TESTCASE}"/multi-file-02/file-02.yml \
  "${TESTCASE}"/multi-file-02/file-03.yml \
  "${TESTCASE}"/multi-file-01/file-01.yml \
  "${TESTCASE}"/multi-file-01/file-02.yml \
  "${TESTCASE}"/multi-file-01/file-03.yml \
  "${TESTCASE}"/multi-file-03/file-01.yml \
  "${TESTCASE}"/multi-file-03/file-02.yml \
  "${TESTCASE}"/multi-file-03/file-03.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="312"
spruce merge \
  "${TESTCASE}"/multi-file-03/file-01.yml \
  "${TESTCASE}"/multi-file-03/file-02.yml \
  "${TESTCASE}"/multi-file-03/file-03.yml \
  "${TESTCASE}"/multi-file-01/file-01.yml \
  "${TESTCASE}"/multi-file-01/file-02.yml \
  "${TESTCASE}"/multi-file-01/file-03.yml \
  "${TESTCASE}"/multi-file-02/file-01.yml \
  "${TESTCASE}"/multi-file-02/file-02.yml \
  "${TESTCASE}"/multi-file-02/file-03.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="132"
spruce merge \
  "${TESTCASE}"/multi-file-01/file-01.yml \
  "${TESTCASE}"/multi-file-01/file-02.yml \
  "${TESTCASE}"/multi-file-01/file-03.yml \
  "${TESTCASE}"/multi-file-03/file-01.yml \
  "${TESTCASE}"/multi-file-03/file-02.yml \
  "${TESTCASE}"/multi-file-03/file-03.yml \
  "${TESTCASE}"/multi-file-02/file-01.yml \
  "${TESTCASE}"/multi-file-02/file-02.yml \
  "${TESTCASE}"/multi-file-02/file-03.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="231"
spruce merge \
  "${TESTCASE}"/multi-file-02/file-01.yml \
  "${TESTCASE}"/multi-file-02/file-02.yml \
  "${TESTCASE}"/multi-file-02/file-03.yml \
  "${TESTCASE}"/multi-file-03/file-01.yml \
  "${TESTCASE}"/multi-file-03/file-02.yml \
  "${TESTCASE}"/multi-file-03/file-03.yml \
  "${TESTCASE}"/multi-file-01/file-01.yml \
  "${TESTCASE}"/multi-file-01/file-02.yml \
  "${TESTCASE}"/multi-file-01/file-03.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="321"
spruce merge \
  "${TESTCASE}"/multi-file-03/file-01.yml \
  "${TESTCASE}"/multi-file-03/file-02.yml \
  "${TESTCASE}"/multi-file-03/file-03.yml \
  "${TESTCASE}"/multi-file-02/file-01.yml \
  "${TESTCASE}"/multi-file-02/file-02.yml \
  "${TESTCASE}"/multi-file-02/file-03.yml \
  "${TESTCASE}"/multi-file-01/file-01.yml \
  "${TESTCASE}"/multi-file-01/file-02.yml \
  "${TESTCASE}"/multi-file-01/file-03.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

TESTCASE=multi-dir-single-file
PERMUTATION="123"
spruce merge \
  "${TESTCASE}"/single-file-01/file.yml \
  "${TESTCASE}"/single-file-02/file.yml \
  "${TESTCASE}"/single-file-03/file.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="213"
spruce merge \
  "${TESTCASE}"/single-file-02/file.yml \
  "${TESTCASE}"/single-file-01/file.yml \
  "${TESTCASE}"/single-file-03/file.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="312"
spruce merge \
  "${TESTCASE}"/single-file-03/file.yml \
  "${TESTCASE}"/single-file-01/file.yml \
  "${TESTCASE}"/single-file-02/file.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="132"
spruce merge \
  "${TESTCASE}"/single-file-01/file.yml \
  "${TESTCASE}"/single-file-03/file.yml \
  "${TESTCASE}"/single-file-02/file.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="231"
spruce merge \
  "${TESTCASE}"/single-file-02/file.yml \
  "${TESTCASE}"/single-file-03/file.yml \
  "${TESTCASE}"/single-file-01/file.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="321"
spruce merge \
  "${TESTCASE}"/single-file-03/file.yml \
  "${TESTCASE}"/single-file-02/file.yml \
  "${TESTCASE}"/single-file-01/file.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

TESTCASE=multi-file
PERMUTATION="123"
spruce merge \
  "${TESTCASE}"/file-01.yml \
  "${TESTCASE}"/file-02.yml \
  "${TESTCASE}"/file-03.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="213"
spruce merge \
  "${TESTCASE}"/file-02.yml \
  "${TESTCASE}"/file-01.yml \
  "${TESTCASE}"/file-03.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="312"
spruce merge \
  "${TESTCASE}"/file-03.yml \
  "${TESTCASE}"/file-01.yml \
  "${TESTCASE}"/file-02.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="132"
spruce merge \
  "${TESTCASE}"/file-01.yml \
  "${TESTCASE}"/file-03.yml \
  "${TESTCASE}"/file-02.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="231"
spruce merge \
  "${TESTCASE}"/file-02.yml \
  "${TESTCASE}"/file-03.yml \
  "${TESTCASE}"/file-01.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="321"
spruce merge \
  "${TESTCASE}"/file-03.yml \
  "${TESTCASE}"/file-02.yml \
  "${TESTCASE}"/file-01.yml \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

TESTCASE=multi-mixed-dirfile
PERMUTATION="123"
spruce merge \
  "${TESTCASE}"/single-mixed-01/file.yml \
  "${TESTCASE}"/single-mixed-01/file.yaml \
  "${TESTCASE}"/single-mixed-01/file.json \
  "${TESTCASE}"/single-mixed-01/file.ejson \
  "${TESTCASE}"/single-mixed-01.yml \
  "${TESTCASE}"/single-mixed-01.yaml \
  "${TESTCASE}"/single-mixed-01.json \
  "${TESTCASE}"/single-mixed-01.ejson \
  "${TESTCASE}"/single-mixed-02/file.yml \
  "${TESTCASE}"/single-mixed-02/file.yaml \
  "${TESTCASE}"/single-mixed-02/file.json \
  "${TESTCASE}"/single-mixed-02/file.ejson \
  "${TESTCASE}"/single-mixed-02.yml \
  "${TESTCASE}"/single-mixed-02.yaml \
  "${TESTCASE}"/single-mixed-02.json \
  "${TESTCASE}"/single-mixed-02.ejson \
  "${TESTCASE}"/single-mixed-03/file.yml \
  "${TESTCASE}"/single-mixed-03/file.yaml \
  "${TESTCASE}"/single-mixed-03/file.json \
  "${TESTCASE}"/single-mixed-03/file.ejson \
  "${TESTCASE}"/single-mixed-03.yml \
  "${TESTCASE}"/single-mixed-03.yaml \
  "${TESTCASE}"/single-mixed-03.json \
  "${TESTCASE}"/single-mixed-03.ejson \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="213"
spruce merge \
  "${TESTCASE}"/single-mixed-02/file.yml \
  "${TESTCASE}"/single-mixed-02/file.yaml \
  "${TESTCASE}"/single-mixed-02/file.json \
  "${TESTCASE}"/single-mixed-02/file.ejson \
  "${TESTCASE}"/single-mixed-02.yml \
  "${TESTCASE}"/single-mixed-02.yaml \
  "${TESTCASE}"/single-mixed-02.json \
  "${TESTCASE}"/single-mixed-02.ejson \
  "${TESTCASE}"/single-mixed-01/file.yml \
  "${TESTCASE}"/single-mixed-01/file.yaml \
  "${TESTCASE}"/single-mixed-01/file.json \
  "${TESTCASE}"/single-mixed-01/file.ejson \
  "${TESTCASE}"/single-mixed-01.yml \
  "${TESTCASE}"/single-mixed-01.yaml \
  "${TESTCASE}"/single-mixed-01.json \
  "${TESTCASE}"/single-mixed-01.ejson \
  "${TESTCASE}"/single-mixed-03/file.yml \
  "${TESTCASE}"/single-mixed-03/file.yaml \
  "${TESTCASE}"/single-mixed-03/file.json \
  "${TESTCASE}"/single-mixed-03/file.ejson \
  "${TESTCASE}"/single-mixed-03.yml \
  "${TESTCASE}"/single-mixed-03.yaml \
  "${TESTCASE}"/single-mixed-03.json \
  "${TESTCASE}"/single-mixed-03.ejson \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="312"
spruce merge \
  "${TESTCASE}"/single-mixed-03/file.yml \
  "${TESTCASE}"/single-mixed-03/file.yaml \
  "${TESTCASE}"/single-mixed-03/file.json \
  "${TESTCASE}"/single-mixed-03/file.ejson \
  "${TESTCASE}"/single-mixed-03.yml \
  "${TESTCASE}"/single-mixed-03.yaml \
  "${TESTCASE}"/single-mixed-03.json \
  "${TESTCASE}"/single-mixed-03.ejson \
  "${TESTCASE}"/single-mixed-01/file.yml \
  "${TESTCASE}"/single-mixed-01/file.yaml \
  "${TESTCASE}"/single-mixed-01/file.json \
  "${TESTCASE}"/single-mixed-01/file.ejson \
  "${TESTCASE}"/single-mixed-01.yml \
  "${TESTCASE}"/single-mixed-01.yaml \
  "${TESTCASE}"/single-mixed-01.json \
  "${TESTCASE}"/single-mixed-01.ejson \
  "${TESTCASE}"/single-mixed-02/file.yml \
  "${TESTCASE}"/single-mixed-02/file.yaml \
  "${TESTCASE}"/single-mixed-02/file.json \
  "${TESTCASE}"/single-mixed-02/file.ejson \
  "${TESTCASE}"/single-mixed-02.yml \
  "${TESTCASE}"/single-mixed-02.yaml \
  "${TESTCASE}"/single-mixed-02.json \
  "${TESTCASE}"/single-mixed-02.ejson \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="132"
spruce merge \
  "${TESTCASE}"/single-mixed-01/file.yml \
  "${TESTCASE}"/single-mixed-01/file.yaml \
  "${TESTCASE}"/single-mixed-01/file.json \
  "${TESTCASE}"/single-mixed-01/file.ejson \
  "${TESTCASE}"/single-mixed-01.yml \
  "${TESTCASE}"/single-mixed-01.yaml \
  "${TESTCASE}"/single-mixed-01.json \
  "${TESTCASE}"/single-mixed-01.ejson \
  "${TESTCASE}"/single-mixed-03/file.yml \
  "${TESTCASE}"/single-mixed-03/file.yaml \
  "${TESTCASE}"/single-mixed-03/file.json \
  "${TESTCASE}"/single-mixed-03/file.ejson \
  "${TESTCASE}"/single-mixed-03.yml \
  "${TESTCASE}"/single-mixed-03.yaml \
  "${TESTCASE}"/single-mixed-03.json \
  "${TESTCASE}"/single-mixed-03.ejson \
  "${TESTCASE}"/single-mixed-02/file.yml \
  "${TESTCASE}"/single-mixed-02/file.yaml \
  "${TESTCASE}"/single-mixed-02/file.json \
  "${TESTCASE}"/single-mixed-02/file.ejson \
  "${TESTCASE}"/single-mixed-02.yml \
  "${TESTCASE}"/single-mixed-02.yaml \
  "${TESTCASE}"/single-mixed-02.json \
  "${TESTCASE}"/single-mixed-02.ejson \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="231"
spruce merge \
  "${TESTCASE}"/single-mixed-02/file.yml \
  "${TESTCASE}"/single-mixed-02/file.yaml \
  "${TESTCASE}"/single-mixed-02/file.json \
  "${TESTCASE}"/single-mixed-02/file.ejson \
  "${TESTCASE}"/single-mixed-02.yml \
  "${TESTCASE}"/single-mixed-02.yaml \
  "${TESTCASE}"/single-mixed-02.json \
  "${TESTCASE}"/single-mixed-02.ejson \
  "${TESTCASE}"/single-mixed-03/file.yml \
  "${TESTCASE}"/single-mixed-03/file.yaml \
  "${TESTCASE}"/single-mixed-03/file.json \
  "${TESTCASE}"/single-mixed-03/file.ejson \
  "${TESTCASE}"/single-mixed-03.yml \
  "${TESTCASE}"/single-mixed-03.yaml \
  "${TESTCASE}"/single-mixed-03.json \
  "${TESTCASE}"/single-mixed-03.ejson \
  "${TESTCASE}"/single-mixed-01/file.yml \
  "${TESTCASE}"/single-mixed-01/file.yaml \
  "${TESTCASE}"/single-mixed-01/file.json \
  "${TESTCASE}"/single-mixed-01/file.ejson \
  "${TESTCASE}"/single-mixed-01.yml \
  "${TESTCASE}"/single-mixed-01.yaml \
  "${TESTCASE}"/single-mixed-01.json \
  "${TESTCASE}"/single-mixed-01.ejson \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="321"
spruce merge \
  "${TESTCASE}"/single-mixed-03/file.yml \
  "${TESTCASE}"/single-mixed-03/file.yaml \
  "${TESTCASE}"/single-mixed-03/file.json \
  "${TESTCASE}"/single-mixed-03/file.ejson \
  "${TESTCASE}"/single-mixed-03.yml \
  "${TESTCASE}"/single-mixed-03.yaml \
  "${TESTCASE}"/single-mixed-03.json \
  "${TESTCASE}"/single-mixed-03.ejson \
  "${TESTCASE}"/single-mixed-02/file.yml \
  "${TESTCASE}"/single-mixed-02/file.yaml \
  "${TESTCASE}"/single-mixed-02/file.json \
  "${TESTCASE}"/single-mixed-02/file.ejson \
  "${TESTCASE}"/single-mixed-02.yml \
  "${TESTCASE}"/single-mixed-02.yaml \
  "${TESTCASE}"/single-mixed-02.json \
  "${TESTCASE}"/single-mixed-02.ejson \
  "${TESTCASE}"/single-mixed-01/file.yml \
  "${TESTCASE}"/single-mixed-01/file.yaml \
  "${TESTCASE}"/single-mixed-01/file.json \
  "${TESTCASE}"/single-mixed-01/file.ejson \
  "${TESTCASE}"/single-mixed-01.yml \
  "${TESTCASE}"/single-mixed-01.yaml \
  "${TESTCASE}"/single-mixed-01.json \
  "${TESTCASE}"/single-mixed-01.ejson \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

TESTCASE=multi-mixed
PERMUTATION="123"
spruce merge \
  "${TESTCASE}"/file-01.yml \
  "${TESTCASE}"/file-01.yaml \
  "${TESTCASE}"/file-01.json \
  "${TESTCASE}"/file-01.ejson \
  "${TESTCASE}"/file-02.yml \
  "${TESTCASE}"/file-02.yaml \
  "${TESTCASE}"/file-02.json \
  "${TESTCASE}"/file-02.ejson \
  "${TESTCASE}"/file-03.yml \
  "${TESTCASE}"/file-03.yaml \
  "${TESTCASE}"/file-03.json \
  "${TESTCASE}"/file-03.ejson \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="213"
spruce merge \
  "${TESTCASE}"/file-02.yml \
  "${TESTCASE}"/file-02.yaml \
  "${TESTCASE}"/file-02.json \
  "${TESTCASE}"/file-02.ejson \
  "${TESTCASE}"/file-01.yml \
  "${TESTCASE}"/file-01.yaml \
  "${TESTCASE}"/file-01.json \
  "${TESTCASE}"/file-01.ejson \
  "${TESTCASE}"/file-03.yml \
  "${TESTCASE}"/file-03.yaml \
  "${TESTCASE}"/file-03.json \
  "${TESTCASE}"/file-03.ejson \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="312"
spruce merge \
  "${TESTCASE}"/file-03.yml \
  "${TESTCASE}"/file-03.yaml \
  "${TESTCASE}"/file-03.json \
  "${TESTCASE}"/file-03.ejson \
  "${TESTCASE}"/file-01.yml \
  "${TESTCASE}"/file-01.yaml \
  "${TESTCASE}"/file-01.json \
  "${TESTCASE}"/file-01.ejson \
  "${TESTCASE}"/file-02.yml \
  "${TESTCASE}"/file-02.yaml \
  "${TESTCASE}"/file-02.json \
  "${TESTCASE}"/file-02.ejson \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="132"
spruce merge \
  "${TESTCASE}"/file-01.yml \
  "${TESTCASE}"/file-01.yaml \
  "${TESTCASE}"/file-01.json \
  "${TESTCASE}"/file-01.ejson \
  "${TESTCASE}"/file-03.yml \
  "${TESTCASE}"/file-03.yaml \
  "${TESTCASE}"/file-03.json \
  "${TESTCASE}"/file-03.ejson \
  "${TESTCASE}"/file-02.yml \
  "${TESTCASE}"/file-02.yaml \
  "${TESTCASE}"/file-02.json \
  "${TESTCASE}"/file-02.ejson \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="231"
spruce merge \
  "${TESTCASE}"/file-02.yml \
  "${TESTCASE}"/file-02.yaml \
  "${TESTCASE}"/file-02.json \
  "${TESTCASE}"/file-02.ejson \
  "${TESTCASE}"/file-03.yml \
  "${TESTCASE}"/file-03.yaml \
  "${TESTCASE}"/file-03.json \
  "${TESTCASE}"/file-03.ejson \
  "${TESTCASE}"/file-01.yml \
  "${TESTCASE}"/file-01.yaml \
  "${TESTCASE}"/file-01.json \
  "${TESTCASE}"/file-01.ejson \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

PERMUTATION="321"
spruce merge \
  "${TESTCASE}"/file-03.yml \
  "${TESTCASE}"/file-03.yaml \
  "${TESTCASE}"/file-03.json \
  "${TESTCASE}"/file-03.ejson \
  "${TESTCASE}"/file-02.yml \
  "${TESTCASE}"/file-02.yaml \
  "${TESTCASE}"/file-02.json \
  "${TESTCASE}"/file-02.ejson \
  "${TESTCASE}"/file-01.yml \
  "${TESTCASE}"/file-01.yaml \
  "${TESTCASE}"/file-01.json \
  "${TESTCASE}"/file-01.ejson \
  > "expected/${TESTCASE}-${PERMUTATION}.yml"

TESTCASE=single-dir-multi-dir-multi-file
spruce merge \
  "${TESTCASE}"/multi-dir-multi-file/multi-file-01/file-01.yml \
  "${TESTCASE}"/multi-dir-multi-file/multi-file-01/file-02.yml \
  "${TESTCASE}"/multi-dir-multi-file/multi-file-01/file-03.yml \
  "${TESTCASE}"/multi-dir-multi-file/multi-file-02/file-01.yml \
  "${TESTCASE}"/multi-dir-multi-file/multi-file-02/file-02.yml \
  "${TESTCASE}"/multi-dir-multi-file/multi-file-02/file-03.yml \
  "${TESTCASE}"/multi-dir-multi-file/multi-file-03/file-01.yml \
  "${TESTCASE}"/multi-dir-multi-file/multi-file-03/file-02.yml \
  "${TESTCASE}"/multi-dir-multi-file/multi-file-03/file-03.yml \
  > "expected/${TESTCASE}.yml"

TESTCASE=single-dir-multi-dir-single-file
spruce merge \
  "${TESTCASE}"/multi-dir-single-file/single-file-01/file.yml \
  "${TESTCASE}"/multi-dir-single-file/single-file-02/file.yml \
  "${TESTCASE}"/multi-dir-single-file/single-file-03/file.yml \
  > "expected/${TESTCASE}.yml"

TESTCASE=single-dir-multi-file
spruce merge \
  "${TESTCASE}"/multi-file/file-01.yml \
  "${TESTCASE}"/multi-file/file-02.yml \
  "${TESTCASE}"/multi-file/file-03.yml \
  > "expected/${TESTCASE}.yml"

TESTCASE=single-dir-multi-mixed-dirfile
spruce merge \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-01/file.yml \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-01/file.yaml \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-01/file.json \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-01/file.ejson \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-01.yml \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-01.yaml \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-01.json \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-01.ejson \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-02/file.yml \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-02/file.yaml \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-02/file.json \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-02/file.ejson \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-02.yml \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-02.yaml \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-02.json \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-02.ejson \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-03/file.yml \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-03/file.yaml \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-03/file.json \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-03/file.ejson \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-03.yml \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-03.yaml \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-03.json \
  "${TESTCASE}"/multi-mixed-dirfile/single-mixed-03.ejson \
  > "expected/${TESTCASE}.yml"

TESTCASE=single-dir-multi-mixed
spruce merge \
  "${TESTCASE}"/multi-mixed/file-01.yml \
  "${TESTCASE}"/multi-mixed/file-01.yaml \
  "${TESTCASE}"/multi-mixed/file-01.json \
  "${TESTCASE}"/multi-mixed/file-01.ejson \
  "${TESTCASE}"/multi-mixed/file-02.yml \
  "${TESTCASE}"/multi-mixed/file-02.yaml \
  "${TESTCASE}"/multi-mixed/file-02.json \
  "${TESTCASE}"/multi-mixed/file-02.ejson \
  "${TESTCASE}"/multi-mixed/file-03.yml \
  "${TESTCASE}"/multi-mixed/file-03.yaml \
  "${TESTCASE}"/multi-mixed/file-03.json \
  "${TESTCASE}"/multi-mixed/file-03.ejson \
  > "expected/${TESTCASE}.yml"

TESTCASE=single-dir-single-file
spruce merge \
  "${TESTCASE}"/single-file/file.yml \
  > "expected/${TESTCASE}.yml"

TESTCASE=single-dir-single-mixed-dirfile
spruce merge \
  "${TESTCASE}"/single-mixed-dirfile/single-mixed/file.yml \
  "${TESTCASE}"/single-mixed-dirfile/single-mixed/file.yaml \
  "${TESTCASE}"/single-mixed-dirfile/single-mixed/file.json \
  "${TESTCASE}"/single-mixed-dirfile/single-mixed/file.ejson \
  "${TESTCASE}"/single-mixed-dirfile/single-mixed.yml \
  "${TESTCASE}"/single-mixed-dirfile/single-mixed.yaml \
  "${TESTCASE}"/single-mixed-dirfile/single-mixed.json \
  "${TESTCASE}"/single-mixed-dirfile/single-mixed.ejson \
  > "expected/${TESTCASE}.yml"

TESTCASE=single-dir-single-mixed
spruce merge \
  "${TESTCASE}"/single-mixed/file.yml \
  "${TESTCASE}"/single-mixed/file.yaml \
  "${TESTCASE}"/single-mixed/file.json \
  "${TESTCASE}"/single-mixed/file.ejson \
  > "expected/${TESTCASE}.yml"

TESTCASE=single-file
spruce merge \
  "${TESTCASE}"/file.yml \
  > "expected/${TESTCASE}.yml"

TESTCASE=single-mixed-dirfile
spruce merge \
  "${TESTCASE}"/single-mixed/file.yml \
  "${TESTCASE}"/single-mixed/file.yaml \
  "${TESTCASE}"/single-mixed/file.json \
  "${TESTCASE}"/single-mixed/file.ejson \
  "${TESTCASE}"/single-mixed.yml \
  "${TESTCASE}"/single-mixed.yaml \
  "${TESTCASE}"/single-mixed.json \
  "${TESTCASE}"/single-mixed.ejson \
  > "expected/${TESTCASE}.yml"

TESTCASE=single-mixed
spruce merge \
  "${TESTCASE}"/file.yml \
  "${TESTCASE}"/file.yaml \
  "${TESTCASE}"/file.json \
  "${TESTCASE}"/file.ejson \
  > "expected/${TESTCASE}.yml"