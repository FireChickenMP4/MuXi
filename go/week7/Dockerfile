FROM deep-inside-the-container:C
# TODO 设置环境变量MUXI=MUXI,将/app/secret/muxi.txt移动到/muxi/muxi.txt      
# journal 11月27日 ：
# 要离开了，如果有人看到这里一定会以为我是个疯子或者是瞥见了古神而丧失了理智 ，其实我也只是一个可怜虫罢了，完成这个dockerfile，把生成的镜像导出成tar，使用form发送到这个地址，参数是image：
# https://deep-inside-the-container.muxixyz.com/finish
# 在那里，我会告诉你一切。
ENV MUXI=MUXI
RUN mkdir -p /muxi && \
    mv /app/secret/muxi.txt /muxi/muxi.txt
CMD ["sh"]