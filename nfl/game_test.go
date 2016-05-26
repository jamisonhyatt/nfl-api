package nfl

import (
    "testing"
    "fmt"
)

func TestRight (t *testing.T) {

    expected := "BC";
    actual := Right("ABC",2);
    if actual != expected {
        t.Error(fmt.Sprintf("Test failed expected %v received %v",expected, actual))
    }

}

func TestRightEqual (t *testing.T) {

    expected := "ABC";
    actual := Right("ABC",3);
    if actual != expected {
        t.Error(fmt.Sprintf("Test failed expected %v received %v",expected, actual))
    }

}

func TestRightLess (t *testing.T) {

    expected := "ABC";
    actual := Right("ABC",5);
    if actual != expected {
        t.Error(fmt.Sprintf("Test failed expected %v received %v",expected, actual))
    }

}