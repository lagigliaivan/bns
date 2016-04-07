package ar.com.bestprice.buyitnow;

import android.util.Log;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.ProtocolException;
import java.net.URL;
import java.util.concurrent.Callable;

/**
 * Created by ivan on 07/04/16.
 */
public class ServiceClient implements Callable {


    private final String URL;

    public ServiceClient(String URL) {
        this.URL = URL;
    }

    @Override
    public String call() {

        HttpURLConnection urlConnection = null;
        BufferedReader reader = null;
        //PurchasesContainer p = null;
        String purchases = "";
        try {
            // http://openweathermap.org/API#forecast
            //URL url = new URL("http://10.33.117.120:8080/catalog/products/");
            //URL url = new URL("http://192.168.0.7:8080/catalog/products/");


            //URL url = new URL("http://10.33.117.120:8080/catalog/purchases?groupBy=month");
            URL url = new URL(this.URL);

            // Create the request to OpenWeatherMap, and open the connection
            urlConnection = (HttpURLConnection) url.openConnection();
            urlConnection.setRequestMethod("GET");
            urlConnection.connect();

            // Read the input stream into a String
            InputStream inputStream = urlConnection.getInputStream();
            StringBuffer buffer = new StringBuffer();
            if (inputStream == null) {
                // Nothing to do.
                return "";
            }
            reader = new BufferedReader(new InputStreamReader(inputStream));

            String line;
            while ((line = reader.readLine()) != null) {
                // Since it's JSON, adding a newline isn't necessary (it won't affect parsing)
                // But it does make debugging a *lot* easier if you print out the completed
                // buffer for debugging.
                buffer.append(line + "\n");
            }

            if (buffer.length() == 0) {
                // Stream was empty.  No point in parsing.
                return "";
            }

             purchases = buffer.toString();

        } catch (MalformedURLException e) {
            e.printStackTrace();
        } catch (ProtocolException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        } finally {
            if (urlConnection != null) {
                urlConnection.disconnect();
            }
            if (reader != null) {
                try {
                    reader.close();
                } catch (final IOException e) {
                    Log.e("PlaceholderFragment", "Error closing stream", e);
                }
            }
        }
        return purchases;
    }
}
