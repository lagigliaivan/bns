package ar.com.bestprice.buyitnow;

import android.util.Log;

import com.google.gson.Gson;

import java.io.BufferedWriter;
import java.io.IOException;
import java.io.OutputStream;
import java.io.OutputStreamWriter;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.ArrayList;
import java.util.Date;
import java.util.concurrent.Callable;

import ar.com.bestprice.buyitnow.dto.Purchase;
import ar.com.bestprice.buyitnow.dto.Purchases;

/**
 * Created by ivan on 08/04/16.
 */
public class POSTServiceClient implements Callable {

    private Purchases purchases;
    private String url;

    public POSTServiceClient(String url, Purchases purchases){
        this.url = url;
        this.purchases = purchases;
    }

    @Override
    public Object call() throws Exception {
        HttpURLConnection urlConnection = null;
        int responseCode = 203;
        // BufferedReader reader = null;
        try {

            URL url = new URL(this.url);

            urlConnection = (HttpURLConnection) url.openConnection();
            urlConnection.setRequestMethod("POST");
            urlConnection.setRequestProperty("Authorization", Context.getContext().getSha1());
            urlConnection.setDoInput(true);
            urlConnection.setDoOutput(true);

            urlConnection.connect();

            Gson gson = new Gson();
            String it = gson.toJson(purchases);

            OutputStream os = urlConnection.getOutputStream();
            BufferedWriter writer = new BufferedWriter(
                    new OutputStreamWriter(os, "UTF-8"));
            //writer.write(getPostDataString(postDataParams));
            writer.write(it);
            writer.flush();
            writer.close();
            os.close();


            responseCode = urlConnection.getResponseCode();


        } catch (IOException e) {
            Log.e("PlaceholderFragment", "Error ", e);
            // If the code didn't successfully get the weather data, there's no point in attemping
            // to parse it.

        } catch (Exception ex) {
            Log.e("PlaceholderFragment", "Error ", ex);
        }
        finally{
            if (urlConnection != null) {
                urlConnection.disconnect();
            }
        }
        return responseCode;
    }
}
