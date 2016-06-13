package ar.com.bestprice.buyitnow;

/**
 * Created by ivan on 19/04/16.
 */
public class Context {

    private static Context context = new Context();
    private String user;
    private String pass;
    private String login;
    private String serviceURL = "http://192.168.0.3:8080/catalog/";

    private Context(){}

    public static Context getContext(){
        return context;
    }

    public String getServiceURL() {
        return serviceURL;
    }

    public String getUser(){
        return user;
    }

    public String getPass(){
        return pass;
    }

    public void setServiceURL(String URL) {
        this.serviceURL = URL;
    }

    public void setUser(String user) {
        this.user = user;
    }

    public void setPass(String pass) {
        this.pass = pass;
    }

    public void setLogin(String login) { this.login = login;}
    public String getLogin() { return this.login;}
}
