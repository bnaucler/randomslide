bin/build.sh
rm data/rs.db
bin/rsserver -v &
sleep 1
curl 'localhost:6291/addtext?tags=apa&ttext=apa1&btext=apa1'
curl 'localhost:6291/addtext?tags=apa&ttext=apa2'
curl 'localhost:6291/addtext?tags=kaka&ttext=kaka1&btext=kaka1'
curl 'localhost:6291/addtext?tags=kaka&ttext=kaka2'
curl 'localhost:6291/gettags'
curl 'localhost:6291/addtext?tags=apa%20kaka%20satan&ttext=all1&btext=all1'
curl 'localhost:6291/gettags'
kill `cat data/rsserver.pid`
rm data/rsserver.pid
bin/dbdump data/rs.db
