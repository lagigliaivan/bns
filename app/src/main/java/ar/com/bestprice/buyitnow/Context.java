package ar.com.bestprice.buyitnow;

import java.io.UnsupportedEncodingException;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Formatter;

/**
 * Created by ivan on 19/04/16.
 */
public class Context {

    private static Context context = new Context();
    private String user = "mayname@gmail.com.ar";
    private String pass = "password";
    private String login;
    private String serviceURL = "http://192.168.0.4:8080/catalog";

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

    public String getSha1()
    {

        String userLogin[] = getUser().split("@");
        String user = userLogin[0];
        String mail = userLogin[1];

        String sha1 = "";
        try
        {
            MessageDigest crypt = MessageDigest.getInstance("SHA-1");
            crypt.reset();
            String credentials = user + ":" + getPass() + "@" + mail ;
            crypt.update(credentials.getBytes("UTF-8"));
            sha1 = byteToHex(crypt.digest());
        }
        catch(NoSuchAlgorithmException e)
        {
            e.printStackTrace();
        }
        catch(UnsupportedEncodingException e)
        {
            e.printStackTrace();
        }
        return sha1;
    }
    private static String byteToHex(final byte[] hash)
    {
        Formatter formatter = new Formatter();
        for (byte b : hash)
        {
            formatter.format("%02x", b);
        }
        String result = formatter.toString();
        formatter.close();
        return result;
    }
}
