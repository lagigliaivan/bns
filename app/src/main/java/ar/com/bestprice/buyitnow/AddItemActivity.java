package ar.com.bestprice.buyitnow;

import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;

/**
 * Created by ivan on 01/04/16.
 */
public class AddItemActivity extends AppCompatActivity{
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        setContentView(R.layout.activity_add_item);

        String barcodeId = getIntent().getStringExtra("BarCode");

        Log.d("Received", barcodeId);
    }
}
