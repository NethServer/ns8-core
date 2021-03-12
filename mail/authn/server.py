import cherrypy

class AuthPassGenerator(object):
    @cherrypy.expose
    def index(self):
        cherrypy.response.headers['Content-Type'] = "text/plain;charset=utf-8"
        cherrypy.response.headers['Auth-Status'] = "OK"
        cherrypy.response.headers['Auth-Server'] = "10.133.0.3"
        cherrypy.response.headers['Auth-Port'] = "25"
        return "auth pass\n"

if __name__ == '__main__':
    cherrypy.config.update({'server.socket_host': '0.0.0.0',
                            'server.socket_port': 9000,
                            })
    cherrypy.quickstart(AuthPassGenerator())