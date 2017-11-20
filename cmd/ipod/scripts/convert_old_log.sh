grep -E  'Input|Output' | perl -pne 's/.*Input.*\[(.*)\]/<$1/g' | perl -pne 's/.*Output.*\[(.*)\]/>$1/g'
