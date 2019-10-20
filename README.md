# RainflowGo
Simple CLI tool to perform rainflow counting to ASTM E1049 85 written in Go

## Usage
An executable file for Linux64 and Win64 has been provided in the dist directory. To use RainflowGo, you can run the executable with 3 flags

### Required
`-i` Path to the input text file, this should be the raw values of the series seperated by new lines, no commas. 

`-r` Size of the range over which the rainflow counting will be grouped, for examlple chosing 2 means counts will be grouped by bins 0-2, 2-4, 4-6 ... etc

### Optional
`-o` Path to output file, this file will be created or overwritten and the results of the analysis will be written as a comma seperated table. If this flag is not provided no output file will be written. 

## Example
This example uses data from Figure 6 of ASTM E1049 85. 

1. First create a text file called `test_data.txt` with stress ranges in the format shown below 
    ```
    -2
    1
    -3
    5
    -1
    3
    4
    4
    -2
    ```
1. Using the command prompt navigate to the directory where the *RainflowGo* executable is located
1. Run the following commands. A range of 1 is specified for the bins and the result will be written to `output.csv`
    1. for Linux

        ```
        ./RainflowGo -i test_data.txt -r 1 -o output.csv
        ```
    1. Or Windows
        ```
        RainflowGo.exe -i test_data.txt -r 1 -o output.csv
        ```
1. In the console this result will be shown. 

    ![](/img/output.png)
    
    The bin mean is the midpoint between the bin low and high values, where as the range mean is the count weighted sum of the stresses found within the bin. 
