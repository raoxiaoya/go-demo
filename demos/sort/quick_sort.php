<?php
function swap(&$a, &$b)
{
    $temp = $a;
    $a = $b;
    $b = $temp;
}
function partition(&$arr, $low, $high)
{
    // 10, 80, 30, 90, 40
    $pivot = $arr[$high];
    $i = $low - 1; // -1

    for ($j = $low; $j <= $high - 1; $j++) {
        if ($arr[$j] < $pivot) {
            $i++;

            swap($arr[$i], $arr[$j]);
        }
    }

    swap($arr[$i + 1], $arr[$high]);
    return $i + 1;
}
function quickSort(&$arr, $low, $high)
{
    if ($low < $high) {
        $pivot = $arr[$high];
        $i = $low - 1;

        for ($j = $low; $j <= $high - 1; $j++) {
            if ($arr[$j] < $pivot) {
                $i++;

                $temp = $arr[$i];
                $arr[$i] = $arr[$j];
                $arr[$j] = $temp;
            }
        }

        $temp = $arr[$i + 1];
        $arr[$i + 1] = $arr[$high];
        $arr[$high] = $temp;

        $pi = $i + 1;

        quickSort($arr, $low, $pi - 1);
        quickSort($arr, $pi + 1, $high);
    }
}

$arr = [];
$len = 10;
for ($i = 0; $i < $len; $i++) {
    $arr[$i] = rand(10, 100000);
}

$start = microtime(true);

// quickSort($arr, 0, $len - 1);

// insertSort($arr);

$res = combineSort($arr);
echo join(',', $res).PHP_EOL;

$end = microtime(true);

echo $end - $start;
echo PHP_EOL;


function insertSort($arr)
{
    $res = [];
    foreach ($arr as $v) {
        if (empty($res)) {
            $res[0] = $v;
        } else {
            if ($v <= $res[0]) {
                array_splice($res, 0, 0, $v);
                continue;
            }
            if ($v > $res[count($res) - 1]) {
                array_push($res, $v);
                continue;
            }
            foreach ($res as $key => $val) {
                if ($v <= $val) {
                    array_splice($res, $key, 0, $v);
                    break;
                }
            }
        }
    }

    echo join(',', $res).PHP_EOL;

    return $res;
}

function combineSort($arr)
{
    $len = count($arr);
    if ($len <= 1) {
        return $arr;
    }

    $mid   = floor($len / 2);
    $left  = array_slice($arr, 0, $mid);
    $right = array_slice($arr, $mid);

    // 递归拆分
    $arr1 = combineSort($left);
    $arr2 = combineSort($right);

    // 合并并排序，从小到大
    $arr3 = [];
    while (!empty($arr1) && !empty($arr2)) {
        array_push(
            $arr3,
            $arr1[0] <= $arr2[0] ? array_shift($arr1) : array_shift($arr2)
        );
    }

    return array_merge($arr3, $arr1, $arr2);// $arr1 和 $arr2 至少有一个为空
}
