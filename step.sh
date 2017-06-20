#!/bin/bash
set -e

# Required parameters
if [ -z "${input_file}" ] ; then
  echo " [!] Missing required input: input_file"
  exit 1
fi

if [ ! -f "${input_file}" ] ; then
  echo " [!] File doesn't exist at specified path: ${input_file}"
  exit 1
fi

if [ -z "${output_file}" ] ; then
  echo " [!] Missing required input: output_file"
  exit 1
fi

echo "This is the value specified for the input 'input_file': ${input_file}"
echo "This is the value specified for the input 'output_file': ${output_file}"
echo "This is the value specified for the input 'move_action': ${move_action}"

if [ "${move_action}" = "move" ] ; then
    mv ${input_file} ${output_file}
else
    cp ${input_file} ${output_file}
fi

if [ "$?" -ne "0" ]; then
  echo "Submission failed"
  exit 1
fi

echo "Submission successful."




#
# --- Export Environment Variables for other Steps:
# You can export Environment Variables for other Steps with
#  envman, which is automatically installed by `bitrise setup`.
# A very simple example:
#envman add --key EXAMPLE_STEP_OUTPUT --value 'the value you want to share'
# Envman can handle piped inputs, which is useful if the text you want to
# share is complex and you don't want to deal with proper bash escaping:
#  cat file_with_complex_input | envman add --KEY EXAMPLE_STEP_OUTPUT
# You can find more usage examples on envman's GitHub page
#  at: https://github.com/bitrise-io/envman

#
# --- Exit codes:
# The exit code of your Step is very important. If you return
#  with a 0 exit code `bitrise` will register your Step as "successful".
# Any non zero exit code will be registered as "failed" by `bitrise`.
