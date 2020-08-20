FROM scratch
COPY main .
COPY YOU-KNOW-be4a1d88e2c3.json .
COPY .env .

#バイナリ実行
CMD ["./main"]