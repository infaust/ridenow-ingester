DATABASE="forecasts"
PG_SUPERUSER="NOSUPERUSER"
USER="ridenow_user"
PASSWORD=$RIDENOW_DB_PASSWORD


echo "Creating db '${DATABASE}' for user '${USER}'"
gosu postgres psql -c "CREATE USER ${USER} CREATEDB ${PG_SUPERUSER} NOCREATEROLE INHERIT LOGIN UNENCRYPTED PASSWORD '${PASSWORD}';"
gosu postgres createdb --owner ${USER} --template template0 --encoding=UTF8 --lc-ctype=en_US.UTF-8 --lc-collate=en_US.UTF-8 ${DATABASE}